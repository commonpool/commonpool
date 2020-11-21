import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {map, pluck, shareReplay, startWith, switchMap} from 'rxjs/operators';
import {GetUserMembershipsRequest, GetUserMembershipsResponse, Membership, UserInfoResponse} from '../../api/models';
import {AuthService} from '../../auth.service';
import {combineLatest, Observable, OperatorFunction, Subject} from 'rxjs';
import {UserInfoService} from '../user-info.service';

@Component({
  selector: 'app-user-groups-view',
  templateUrl: './user-groups-view.component.html',
  styleUrls: ['./user-groups-view.component.css']
})
export class UserGroupsViewComponent implements OnInit {

  refreshSubject: Subject<any>;
  userInfo$: Observable<UserInfoResponse>;
  username$: Observable<string>;
  memberships$: Observable<Membership[]>;
  activeMembership$: Observable<Membership[]>;
  incomingMembership$: Observable<Membership[]>;
  outgoingMembership$: Observable<Membership[]>;
  isMyProfile$: Observable<boolean>;

  constructor(
    private route: ActivatedRoute,
    private backend: BackendService,
    private auth: AuthService,
    private userInfoService: UserInfoService
  ) {
  }

  ngOnInit(): void {
    this.userInfo$ = this.userInfoService.getUserInfo();
    this.username$ = this.userInfo$.pipe(pluck('username'));
    this.refreshSubject = new Subject();
    this.memberships$ = combineLatest([
      this.refreshSubject.asObservable().pipe(startWith(true)),
      this.userInfoService.getUserId()
    ]).pipe(
      map(([_, userId]) => userId),
      getUserMemberships(this.backend),
      pluck('memberships'),
      shareReplay()
    );
    this.activeMembership$ = this.memberships$.pipe(filterActiveMemberships());
    this.incomingMembership$ = this.memberships$.pipe(filterIncomingMemberships());
    this.outgoingMembership$ = this.memberships$.pipe(filterOutgoingMemberships());
    this.isMyProfile$ = this.userInfoService.getIsMyProfile();
  }

  refresh() {
    this.refreshSubject.next();
  }
}

function filterActiveMemberships(): OperatorFunction<Membership[], Membership[]> {
  return source => {
    return source.pipe(map(ms => ms.filter(m => m.groupConfirmed && m.userConfirmed)));
  };
}

function filterIncomingMemberships(): OperatorFunction<Membership[], Membership[]> {
  return source => {
    return source.pipe(map(ms => ms.filter(m => m.groupConfirmed && !m.userConfirmed)));
  };
}

function filterOutgoingMemberships(): OperatorFunction<Membership[], Membership[]> {
  return source => {
    return source.pipe(map(ms => ms.filter(m => !m.groupConfirmed && m.userConfirmed)));
  };
}

function getUserMemberships(backend: BackendService): OperatorFunction<string, GetUserMembershipsResponse> {
  return source => {
    return source.pipe(switchMap(id => backend.getUserMemberships(new GetUserMembershipsRequest(id))));
  };
}
