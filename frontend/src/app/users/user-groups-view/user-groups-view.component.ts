import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {map, pluck, shareReplay, startWith, switchMap, tap} from 'rxjs/operators';
import {GetUserMembershipsRequest, MembershipStatus} from '../../api/models';
import {AuthService} from '../../auth.service';
import {combineLatest, Subject} from 'rxjs';

@Component({
  selector: 'app-user-groups-view',
  templateUrl: './user-groups-view.component.html',
  styleUrls: ['./user-groups-view.component.css']
})
export class UserGroupsViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService, private auth: AuthService) {
  }

  refreshSubject = new Subject();
  refresh$ = this.refreshSubject.asObservable().pipe(startWith(true));
  userId$ = this.route.parent.params.pipe(pluck('id'));

  memberships$ = combineLatest([this.refresh$, this.userId$]).pipe(
    map(([_, userId]) => userId),
    tap(console.log),
    switchMap(id => this.backend.getUserMemberships(new GetUserMembershipsRequest(id))),
    pluck('memberships'),
    shareReplay()
  );

  activeMembership$ = this.memberships$.pipe(map(ms => ms.filter(m => m.groupConfirmed && m.userConfirmed)));
  incomingMembership$ = this.memberships$.pipe(map(ms => ms.filter(m => m.groupConfirmed && !m.userConfirmed)));
  outgoingMembership$ = this.memberships$.pipe(map(ms => ms.filter(m => !m.groupConfirmed && m.userConfirmed)));

  authUser$ = this.auth.session$.pipe(pluck('id'));
  isMyProfile$ = combineLatest([this.userId$, this.authUser$]).pipe(
    map(([userId, authUserId]) => userId === authUserId),
    shareReplay()
  );

  ngOnInit(): void {
  }

}
