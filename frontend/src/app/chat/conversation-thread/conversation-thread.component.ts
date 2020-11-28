import {
  AfterViewInit,
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  QueryList,
  ViewChild,
  ViewChildren
} from '@angular/core';
import {BehaviorSubject, combineLatest, Observable, Subject} from 'rxjs';
import {filter, map, pluck, startWith, switchMap} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {Event, EventType, GetMessagesResponse, Message, SectionBlock, TextObject, TextType} from '../../api/models';
import {format} from 'date-fns';
import {ChatService, ConversationType} from '../chat.service';
import {MapMessages} from '../utils/messages-mapper';

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
  styles: [
      `
      .input-section {
        position: absolute;
        bottom: 0;
        left: 0;
        height: 2.5rem;
        z-index: 1000;
        width: 100%
      }

      textarea {
        height: 2.5rem;
        resize: none;
      }

      .messages-container {
        height: calc(100vh - 11rem);
        overflow-y: auto;
      }

      .messages-inner-container {
        max-width: 60rem;
      }
    `
  ],
  template: `
    <div class="channel w-100 " style="position:relative">

      <div style="height: 3rem; position: relative; top:0; left: 0; width: 100%; z-index:1000"
           class="border-bottom bg-primary">
        <a class="btn btn-light mt-1 ml-1" [routerLink]="'/messages'">
          🡠
        </a>
      </div>

      <div class="py-2 w-100 messages-container" #scrollframe (scroll)="scrolled($event)" (resize)="scrolled($event)">
        <div class="messages messages-inner-container" style="line-break: loose">
          <ng-container *ngIf="messageGroups$ | async; let messageGroups">
            <ng-container *ngFor="let messageGroup of messageGroups;">
              <app-message-group [messageGroup]="messageGroup" #item></app-message-group>
            </ng-container>
          </ng-container>
        </div>
      </div>

      <div class="input-section px-2 mb-3 d-flex flex-row">
        <textarea class="form-control rounded-0"
                  [(ngModel)]="content"
                  (keydown.enter)="sendMessage($event)">
        </textarea>
        <div class="send">
          <button class="btn btn-primary ml-2" (click)="sendMessage($event)">
            <app-arrow-right></app-arrow-right>
          </button>
        </div>
      </div>

    </div>

  `,
  styleUrls: ['./conversation-thread.component.css']
})
export class ConversationThreadComponent implements OnInit, OnDestroy, AfterViewInit {

  @ViewChild('scrollframe', {static: false}) scrollFrame: ElementRef<HTMLElement>;
  @ViewChildren('item') itemElements: QueryList<any>;
  private scrollContainer: HTMLElement;

  private isNearBottom = true;
  private type: string;
  private skipSubject = new Subject<number>();
  private skip$ = this.skipSubject.pipe(startWith(0));
  private takeSubject = new Subject<number>();
  private take$ = this.takeSubject.pipe(startWith(30));
  private channelIdSubject = new BehaviorSubject<string>(undefined);
  private channelId$ = this.channelIdSubject.asObservable();
  public content = '';
  private routeSub = this.route.params.pipe(pluck('id')).subscribe(channelId => {
    this.chat.setCurrentConversation({
      type: ConversationType.Channel,
      id: channelId
    });
    this.channelIdSubject.next(channelId);
  });

  pub$ = combineLatest<Observable<string>, Observable<Event>>([this.channelIdSubject.asObservable(), this.backend.events$]).pipe(
    filter(([channel, message]) => !!channel && !!message),
    filter(([channel, message]) => message.type === EventType.MessageEvent),
    filter(([channel, message]) => channel === message.channel)
  ).subscribe((a) => {
    this.refresh();
  });

  private triggerSubject = new Subject<void>();
  private trigger$ = this.triggerSubject.asObservable().pipe(startWith([undefined]));

  public messages$: Observable<GetMessagesResponse> = combineLatest([this.skip$, this.take$, this.channelId$, this.trigger$])
    .pipe(
      switchMap(([s, t, channelId, _]) => {
        return this.backend.getMessages(channelId, new Date().valueOf(), t);
      }),
    );

  public messageGroups$ = this.messages$.pipe(
    pluck<GetMessagesResponse, Message[]>('messages'),
    map((messages: Message[]) => MapMessages(messages)));

  trackMessage(i, o: Message) {
    return o.id;
  }

  constructor(private route: ActivatedRoute, private backend: BackendService, public chat: ChatService) {
    this.type = this.route.snapshot.data.type;
  }

  ngOnInit(): void {

  }

  refresh() {
    this.triggerSubject.next(null);
  }

  sendMessage(event: any) {
    event.preventDefault();
    this.backend.sendMessage(this.channelIdSubject.value, this.content).subscribe(() => {
      this.content = '';
    });
  }

  ngOnDestroy(): void {
    this.chat.setCurrentConversation(undefined);
  }

  private onItemElementsChanged(): void {
    console.log('element changed', this.isNearBottom);
    if (this.isNearBottom) {
      this.scrollToBottom();
    }
  }

  private scrollToBottom(): void {
    this.scrollContainer.scroll({
      top: this.scrollContainer.scrollHeight,
      left: 0,
      behavior: 'auto'
    });
  }

  ngAfterViewInit() {
    this.scrollContainer = this.scrollFrame.nativeElement;
    this.itemElements.changes.subscribe(_ => this.onItemElementsChanged());
  }

  private isUserNearBottom(): boolean {
    const threshold = 300;
    const position = this.scrollContainer.scrollTop + this.scrollContainer.offsetHeight;
    const height = this.scrollContainer.scrollHeight;
    return position > height - threshold;
  }

  scrolled(event: any): void {
    this.isNearBottom = this.isUserNearBottom();
  }

}
