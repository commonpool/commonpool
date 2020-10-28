import {Component, OnInit} from '@angular/core';
import {BehaviorSubject, combineLatest, Subject} from 'rxjs';
import {map, pluck, startWith, switchMap} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {Message} from '../../api/models';
import {format} from 'date-fns';

class MessageGroup {
  constructor(public date: Date, public dateStr: string, public messages: Message[]) {
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
  private trigger$ = this.triggerSubject.asObservable().pipe(startWith(undefined));

  public messages$ = combineLatest([this.skip$, this.take$, this.topic$, this.trigger$]).pipe(switchMap(([skip, take, topic]) => {
    return this.backend.getMessages(topic, skip, take);
  }));

  public messageGroups$ = this.messages$.pipe(
    pluck('messages'),
    map((messages) => {

      const messageGroups: MessageGroup[] = [];
      let lastDate: Date = undefined;
      let lastDateYear: number;
      let lastDateMonth: number;
      let lastDateDay: number;


      for (const message of messages) {
        if (lastDate === undefined
          || lastDate.getFullYear() !== message.sentAtDate.getFullYear()
          || lastDate.getMonth() !== message.sentAtDate.getMonth()
          || lastDate.getDate() !== message.sentAtDate.getDate()) {

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

          const newMessageGroup = new MessageGroup(lastDate, dateStr, []);
          messageGroups.push(newMessageGroup);
        }
        messageGroups[messageGroups.length - 1].messages.push(message);
      }

      return messageGroups;
    })
  );

  trackMessage(i, o: Message) {
    return o.id;
  };

  trackMessageGroup(i, o: MessageGroup) {
    return o.date.toISOString();
  }


  constructor(private route: ActivatedRoute, private backend: BackendService) {
  }

  ngOnInit(): void {
    setInterval(() => this.refresh(), 1000);
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
