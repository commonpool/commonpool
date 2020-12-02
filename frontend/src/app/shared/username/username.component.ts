import {Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {Subject} from 'rxjs';
import {distinctUntilChanged, pluck, shareReplay, switchMap, tap} from 'rxjs/operators';

@Component({
  selector: 'app-username',
  templateUrl: './username.component.html',
  styleUrls: ['./username.component.css']
})
export class UsernameComponent implements OnInit, OnDestroy {

  constructor(private backend: BackendService) {

  }

  private idSubject = new Subject<string>();
  id$ = this.idSubject.asObservable().pipe(
    distinctUntilChanged(),
    shareReplay()
  );

  userInfoSub = this.id$.pipe(
    switchMap(id => this.backend.getUserInfo(id)),
    pluck('username')
  ).subscribe((username) => {
    this.username.next(username);
  });

  ngOnInit(): void {
  }

  @Input()
  set id(value: string) {
    this.idSubject.next(value);
  }

  @Output()
  username: EventEmitter<string> = new EventEmitter<string>();

  ngOnDestroy(): void {
    if (this.userInfoSub) {
      this.userInfoSub.unsubscribe();
    }
  }

}
