import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {Subject} from 'rxjs';
import {pluck, switchMap} from 'rxjs/operators';

@Component({
  selector: 'app-username',
  templateUrl: './username.component.html',
  styleUrls: ['./username.component.css']
})
export class UsernameComponent implements OnInit {

  constructor(private backend: BackendService) {
  }

  private idSubject = new Subject<string>();
  userInfo$ = this.idSubject.asObservable().pipe(
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

}
