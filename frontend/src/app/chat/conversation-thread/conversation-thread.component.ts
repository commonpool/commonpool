import {Component, OnInit} from '@angular/core';
import {combineLatest, Subject} from 'rxjs';
import {pluck, startWith, switchMap} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';

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
  private topicId$ = this.route.params.pipe(pluck('id'));
  public messages$ = combineLatest([this.skip$, this.take$, this.topicId$]).pipe(switchMap(([skip, take, topic]) => {
    return this.backend.getMessages(topic, skip, take);
  }));

  constructor(private route: ActivatedRoute, private backend: BackendService) {
  }

  ngOnInit(): void {
  }

}
