import {Component, OnInit} from '@angular/core';
import {delay, map, pluck, share, shareReplay, startWith, switchMap} from 'rxjs/operators';
import {GetGroupMembershipsRequest, GetMembershipRequest, GetUserMembershipsRequest, MembershipStatus} from '../../api/models';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {AuthService} from '../../auth.service';
import {combineLatest, Subject} from 'rxjs';

@Component({
  selector: 'app-group-invites-view',
  templateUrl: './group-invites-view.component.html',
  styleUrls: ['./group-invites-view.component.css']
})
export class GroupInvitesViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService, private auth: AuthService) {
  }

  refreshSubject = new Subject<boolean>();
  refresh$ = this.refreshSubject.asObservable().pipe(startWith(true), delay(100));
  groupId$ = this.route.parent.params.pipe(pluck('id'));

  memberships$ = combineLatest([this.refresh$, this.groupId$]).pipe(
    map(([_, userId]) => userId),
    switchMap(id => {
      return this.backend.getGroupMemberships(new GetGroupMembershipsRequest(id));
    }),
    pluck('memberships'),
    shareReplay()
  );

  incomingInvitations$ = this.memberships$.pipe(map(memberships => {
    return memberships.filter(m => m.userConfirmed === true && m.groupConfirmed === false);
  }));

  outgoingInvitations$ = this.memberships$.pipe(map(memberships => {
    return memberships.filter(m => m.userConfirmed === false && m.groupConfirmed === true);
  }));

  authUserId$ = this.auth.session$.pipe(
    pluck('id'),
    shareReplay()
  );
  myMembership$ = combineLatest([this.groupId$, this.authUserId$]).pipe(
    switchMap(([groupId, authUserId]) => this.backend.getMembership(new GetMembershipRequest(authUserId, groupId))),
    pluck('membership'),
    shareReplay()
  );
  iAmAdmin$ = this.myMembership$.pipe(
    map(m => m.isAdmin),
    shareReplay()
  );

  ngOnInit(): void {
  }


}
