import {Component, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {debounceTime, filter, map, pluck, shareReplay, startWith, switchMap} from 'rxjs/operators';
import {merge, Observable, Subject} from 'rxjs';
import {ChannelType, EventType, GetMyMembershipsRequest, Membership, Subscription} from '../../api/models';
import {ChatService} from '../chat.service';

@Component({
  selector: 'app-conversation-thread-list',
  templateUrl: './conversation-thread-list.component.html',
  styleUrls: ['./conversation-thread-list.component.css']
})
export class ConversationThreadListComponent implements OnInit {

  constructor(private backend: BackendService, public chat: ChatService) {
  }

  refresh = new Subject();
  refresh$ = this.refresh.asObservable().pipe(startWith([undefined]));
  events$ = this.backend.events$.pipe(
    filter((m) => m.type === EventType.MessageEvent),
  );

  subscriptions$: Observable<Subscription[]> = merge(this.refresh$, this.events$).pipe(
    debounceTime(100),
    switchMap(() => this.backend.getSubscriptions(0, 10).pipe(pluck('subscriptions'))),
    shareReplay()
  );

  conversations$ = this.subscriptions$.pipe(
    map(a => a.filter(s => s.type === ChannelType.Conversation))
  );

  groups$ = this.subscriptions$.pipe(
    map(a => a.filter(s => s.type === ChannelType.Group))
  );

  ngOnInit(): void {

  }

}
