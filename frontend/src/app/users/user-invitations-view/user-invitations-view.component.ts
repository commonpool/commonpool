import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {filter, map, pluck, shareReplay, switchMap} from 'rxjs/operators';
import {GetUserMembershipsRequest, MembershipStatus} from '../../api/models';
import {combineLatest} from 'rxjs';
import {AuthService} from '../../auth.service';

@Component({
  selector: 'app-user-invitations-view',
  templateUrl: './user-invitations-view.component.html',
  styleUrls: ['./user-invitations-view.component.css']
})
export class UserInvitationsViewComponent implements OnInit {


  constructor(private route: ActivatedRoute, private backend: BackendService, private auth: AuthService) {
  }

  userId$ = this.route.parent.params.pipe(pluck('id'));
  memberships$ = this.userId$.pipe(
    switchMap(id => this.backend.getUserMemberships(new GetUserMembershipsRequest(id, MembershipStatus.PendingStatus))),
    pluck('memberships'),
    shareReplay()
  );

  incomingInvitations$ = this.memberships$.pipe(map(memberships => {
    return memberships.filter(m => m.groupConfirmed === true && m.userConfirmed === false);
  }));

  outgoingInvitations$ = this.memberships$.pipe(map(memberships => {
    return memberships.filter(m => m.groupConfirmed === false && m.userConfirmed === true);
  }));

  authUser$ = this.auth.session$.pipe(pluck('id'));
  isMyProfile$ = combineLatest([this.userId$, this.authUser$]).pipe(
    map(([userId, authUserId]) => userId === authUserId),
    shareReplay()
  );

  ngOnInit(): void {
  }

}
