import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {map, pluck, shareReplay, switchMap} from 'rxjs/operators';
import {GetUserMembershipsRequest, MembershipStatus} from '../../api/models';
import {AuthService} from '../../auth.service';
import {combineLatest} from 'rxjs';

@Component({
  selector: 'app-user-groups-view',
  templateUrl: './user-groups-view.component.html',
  styleUrls: ['./user-groups-view.component.css']
})
export class UserGroupsViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService, private auth: AuthService) {
  }

  userId$ = this.route.parent.params.pipe(pluck('id'));
  memberships$ = this.userId$.pipe(
    switchMap(id => this.backend.getUserMemberships(new GetUserMembershipsRequest(id, MembershipStatus.ApprovedMembershipStatus))),
    pluck('memberships'),
    shareReplay()
  );
  authUser$ = this.auth.session$.pipe(pluck('id'));
  isMyProfile$ = combineLatest([this.userId$, this.authUser$]).pipe(
    map(([userId, authUserId]) => userId === authUserId),
    shareReplay()
  );

  ngOnInit(): void {
  }

}
