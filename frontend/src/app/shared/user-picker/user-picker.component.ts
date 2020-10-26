import {Component, forwardRef, Input, OnInit} from '@angular/core';
import {combineLatest, of, Subject} from 'rxjs';
import {map, pluck, startWith, switchMap, tap, withLatestFrom} from 'rxjs/operators';
import {BackendService} from '../../api/backend.service';
import {SearchUsersQuery, UserInfoResponse} from '../../api/models';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';

@Component({
  selector: 'app-user-picker',
  templateUrl: './user-picker.component.html',
  styleUrls: ['./user-picker.component.css'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => UserPickerComponent),
      multi: true
    }
  ]
})
export class UserPickerComponent implements OnInit, ControlValueAccessor {

  // ControlValueAccessor backing field
  private propagateChangeFn: (val: string) => void;

  // Pass in a predicate to filter out UserInfo values
  private predicateSubject = new Subject<(val: string) => boolean>();
  private predicate$ = this.predicateSubject.asObservable().pipe(startWith(() => true));

  @Input()
  set predicate(predicate: (val: string) => boolean) {
    this.predicateSubject.next(predicate);
  }

  // This is the control value

  _selectedUser: string;

  get selectedUser() {
    return this._selectedUser;
  }

  set selectedUser(selectedUser: string) {
    this._selectedUser = selectedUser;
    this.propagateChange(this._selectedUser);
  }

  // Observe ng-select typeahead input and query the backend
  querySubject = new Subject<string>();
  items$ = this.querySubject.asObservable().pipe(
    startWith(''),
    switchMap(q => this.backend.searchUsers(new SearchUsersQuery(q, 10, 0))),
    pluck('users'),
    switchMap((users) => combineLatest([of(users), this.predicate$])),
    tap(([users, predicate]) => {
      if (this.selectedUser && predicate && !predicate(this.selectedUser)) {
        setTimeout(() => {
          this.selectedUser = undefined;
        }, 0);
      }
    }),
    map(([users, predicate]) => users.filter(u => predicate(u.id)))
  );


  constructor(private backend: BackendService) {
  }

  // Begin OnInit implementation

  ngOnInit(): void {

  }

  // End OnInit implementation

  // Begin ControlValueAccessor implementation

  propagateChange(val: string) {
    if (this.propagateChangeFn) {
      this.propagateChangeFn(val);
    }
  }

  registerOnChange(fn: any): void {
    this.propagateChangeFn = fn;
  }

  registerOnTouched(fn: any): void {
  }

  setDisabledState(isDisabled: boolean): void {
  }

  writeValue(obj: any): void {
    this.selectedUser = obj;
  }

  // End ControlValueAccessor implementation


}
