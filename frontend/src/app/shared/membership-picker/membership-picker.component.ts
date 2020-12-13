import {Component, forwardRef, Input, OnInit} from '@angular/core';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';
import {BackendService} from '../../api/backend.service';
import {AuthService} from '../../auth.service';
import {distinctUntilChanged, filter, shareReplay, switchMap, tap} from 'rxjs/operators';
import {GetUserMembershipsRequest, MembershipStatus} from '../../api/models';
import {ReplaySubject} from 'rxjs';

@Component({
  selector: 'app-membership-picker',
  templateUrl: './membership-picker.component.html',
  styleUrls: ['./membership-picker.component.css'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => MembershipPickerComponent),
      multi: true
    }
  ]
})
export class MembershipPickerComponent implements OnInit, ControlValueAccessor {

  private propagateChangeFn: (val: string) => void;

  @Input()
  disabled = false

  constructor(private backend: BackendService, private authService: AuthService) {
  }

  memberships = this.authService.authUserId$.pipe(
    filter(authUserId => !!authUserId),
    switchMap(userId => this.backend.getUserMemberships(new GetUserMembershipsRequest(userId, MembershipStatus.ApprovedMembershipStatus))),
    tap((memberships) => {
      if (memberships.memberships.length === 1) {
        this.selectedMembership.next(memberships.memberships[0].groupId);
      }
    }),
    shareReplay()
  );

  selectedMembership = new ReplaySubject<string>();
  selectedMembership$ = this.selectedMembership.asObservable();

  selectedMembershipSub = this.selectedMembership.pipe(
    distinctUntilChanged()
  ).subscribe(selectedId => {
    if (this.propagateChangeFn) {
      this.propagateChangeFn(selectedId);
    }
  });

  ngOnInit(): void {
  }

  registerOnChange(fn: any): void {
    this.propagateChangeFn = fn;
  }

  registerOnTouched(fn: any): void {
  }

  setDisabledState(isDisabled: boolean): void {
    console.log(isDisabled)
    this.disabled = isDisabled
  }

  writeValue(obj: any): void {
    this.selectedMembership.next(obj);
  }

}
