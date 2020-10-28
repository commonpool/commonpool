import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {LeaveGroupRequest, Membership} from '../../api/models';
import {AuthService} from '../../auth.service';
import {map, pluck, shareReplay, switchMap} from 'rxjs/operators';
import {combineLatest, ReplaySubject, Subject} from 'rxjs';
import {BackendService} from '../../api/backend.service';
import {GroupService} from '../group.service';

@Component({
  selector: 'app-kick-or-leave-group-button',
  templateUrl: './kick-or-leave-group-button.component.html',
  styleUrls: ['./kick-or-leave-group-button.component.css']
})
export class KickOrLeaveGroupButtonComponent implements OnInit {

  constructor(private auth: AuthService, private backend: BackendService, private groupService: GroupService) {
  }

  authUserId$ = this.auth.session$.pipe(pluck('id'));
  membershipSubject = new ReplaySubject<Membership>();
  membership$ = this.membershipSubject.asObservable().pipe(shareReplay());
  isMyMembership$ = combineLatest(([this.authUserId$, this.membership$])).pipe(
    map(([authUserId, membership]) => {
      return membership.userId === authUserId;
    })
  );
  isOwner$ = combineLatest(([this.authUserId$, this.membership$])).pipe(
    map(([authUserId, membership]) => {
      return membership.userId === authUserId && membership.isOwner;
    })
  );
  isAdmin = combineLatest(([this.authUserId$, this.membership$])).pipe(
    map(([authUserId, membership]) => {
      return membership.userId === authUserId && membership.isOwner;
    })
  );
  isActive$ = this.membership$.pipe(map(m => m.userConfirmed && m.groupConfirmed));
  isPending$ = this.membership$.pipe(map(m => !m.userConfirmed || !m.groupConfirmed));
  userConfirmed$ = this.membership$.pipe(map(m => m.userConfirmed));
  groupConfirmed$ = this.membership$.pipe(map(m => m.groupConfirmed));

  myMembership$ = combineLatest([this.authUserId$, this.membership$]).pipe(
    switchMap(([authUserId, membership]) => this.groupService.getPermission(authUserId, membership.groupId))
  );
  iAmAdmin$ = this.myMembership$.pipe(map(m => m.isAdmin));

  pending = false;
  error = undefined;
  @Output()
  left = new EventEmitter<string>();

  kickOut(groupId: string, userId: string) {
    this.pending = true;
    this.error = undefined;
    this.backend.leaveGroup(new LeaveGroupRequest(userId, groupId)).subscribe(res => {
      this.pending = false;
      this.left.next(userId);
    }, err => {
      this.pending = false;
      this.error = err;
    }, () => {
      this.pending = false;
    });
  }


  @Input()
  set membership(value: Membership) {
    this.membershipSubject.next(value);
  }

  ngOnInit(): void {
  }

}
