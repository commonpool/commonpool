import {Component, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {pluck} from 'rxjs/operators';
import {Observable} from 'rxjs';
import {Thread} from '../../api/models';

@Component({
  selector: 'app-conversation-thread-list',
  templateUrl: './conversation-thread-list.component.html',
  styleUrls: ['./conversation-thread-list.component.css']
})
export class ConversationThreadListComponent implements OnInit {

  constructor(private backend: BackendService) {
  }

  threads$: Observable<Thread[]> = this.backend.getThreads(0, 10).pipe(pluck('threads'));

  ngOnInit(): void {

  }

}
