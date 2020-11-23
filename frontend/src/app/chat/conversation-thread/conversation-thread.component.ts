import {Component, OnInit} from '@angular/core';
import {BehaviorSubject, combineLatest, Observable, Subject} from 'rxjs';
import {map, pluck, startWith, switchMap} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {GetMessagesResponse, Message, SectionBlock, TextObject, TextType} from '../../api/models';
import {format} from 'date-fns';

class MessageGroup {
  constructor(public date: Date, public dateStr: string, public messages: Message[]) {
  }
}

enum DisplayType {
  Date = 'date',
  UserMessages = 'userMessages'
}

class DisplayElement {
  constructor(public type: DisplayType) {
  }
}

class DateSeparator extends DisplayElement {
  constructor(public date: Date, public dateStr: string) {
    super(DisplayType.Date);
  }
}

class UserMessages extends DisplayElement {
  constructor(public username: string, public userID: string, public messages: Message[]) {
    super(DisplayType.UserMessages);
  }
}

@Component({
  selector: 'app-conversation-thread',
  templateUrl: './conversation-thread.component.html',
  styleUrls: ['./conversation-thread.component.css']
})
export class ConversationThreadComponent implements OnInit {

  private skipSubject = new Subject<number>();
  private skip$ = this.skipSubject.pipe(startWith(0));
  private takeSubject = new Subject<number>();
  private take$ = this.takeSubject.pipe(startWith(10));
  private topicSubject = new BehaviorSubject<string>(undefined);
  private topic$ = this.topicSubject.asObservable();
  public content = '';
  private topicSub = this.route.params.pipe(pluck('id')).subscribe(topic => {
    this.topicSubject.next(topic);
  });

  private triggerSubject = new Subject<void>();
  private trigger$ = this.triggerSubject.asObservable().pipe(startWith([undefined]));

  public messages$: Observable<GetMessagesResponse> = combineLatest([this.skip$, this.take$, this.topic$, this.trigger$])
    .pipe(
      switchMap(([s, t, topic, _]) => {
        return this.backend.getMessages(topic, s, t);
      }),
    );

  public displayElements$ = this.messages$.pipe(
    pluck<GetMessagesResponse, Message[]>('messages'),
    map((messages: Message[]) => {
      return messages.map(m => {
        if (((!m.blocks || m.blocks.length === 0) && m.text)) {
          m.blocks = [
            new SectionBlock(new TextObject(TextType.PlainTextType, m.text))
          ];
        }
        return m;
      });
    }),
    map((messages: Message[]) => {

      const messageGroups: DisplayElement[] = [];

      let lastDate: Date;
      let lastDateYear: number;
      let lastDateMonth: number;
      let lastDateDay: number;
      let lastUser: string;
      let currentUserMsgGrp: UserMessages;

      for (let i = 0; i < messages.length; i++) {

        const message = messages[i];

        if (i === 0) {
          lastDate = message.sentAtDate;
          lastDateYear = lastDate.getFullYear();
          lastDateMonth = lastDate.getMonth();
          lastDateDay = lastDate.getDate();
        }

        if (this.isDifferentDate(lastDate, message.sentAtDate) || (messages.length - 1 === i)) {

          lastDate = message.sentAtDate;
          lastDateYear = message.sentAtDate.getFullYear();
          lastDateMonth = message.sentAtDate.getMonth();
          lastDateDay = message.sentAtDate.getDate();

          let dateStr = format(lastDate, 'EEEE LLL. Mo.');

          if (this.isToday(lastDate)) {
            dateStr = 'today';
          } else if (this.isYesterday(lastDate)) {
            dateStr = 'yesterday';
          }

          const newDateGroup = new DateSeparator(lastDate, dateStr);
          console.log(newDateGroup);
          messageGroups.push(newDateGroup);

        }

        if (currentUserMsgGrp === undefined || lastUser !== message.sentBy) {
          currentUserMsgGrp = new UserMessages(message.sentByUsername, message.sentBy, []);
          lastUser = message.sentBy;
          messageGroups.push(currentUserMsgGrp);
        }

        currentUserMsgGrp.messages.push(message);
      }

      return messageGroups;
    })
  );

  private isDifferentDate(date1: Date, date2: Date) {
    return (date1 && date1.getFullYear() !== date2.getFullYear())
      || (date1 && date1.getMonth() !== date2.getMonth())
      || (date1 && date1.getDate() !== date2.getDate());

  }

  trackMessage(i, o: Message) {
    return o.id;
  }

  trackMessageGroup(i, o: MessageGroup) {
    return o?.date?.toISOString();
  }

  constructor(private route: ActivatedRoute, private backend: BackendService) {
  }

  ngOnInit(): void {
    setInterval(() => this.refresh(), 5000);
  }

  refresh() {
    this.triggerSubject.next(null);
  }

  sendMessage(event: any) {
    event.preventDefault();
    this.backend.sendMessage(this.topicSubject.value, this.content).subscribe(() => {
      this.content = '';
      this.refresh();
    });
  }

  isToday(someDate) {
    const today = new Date();
    return someDate.getDate() === today.getDate() &&
      someDate.getMonth() === today.getMonth() &&
      someDate.getFullYear() === today.getFullYear();
  }

  isYesterday(someDate) {
    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);
    return someDate.getDate() === yesterday.getDate() &&
      someDate.getMonth() === yesterday.getMonth() &&
      someDate.getFullYear() === yesterday.getFullYear();
  }

}
