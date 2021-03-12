import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {delay, map, pluck, startWith, switchMap, tap} from 'rxjs/operators';
import {BackendService} from '../../api/backend.service';
import {
  GetGroupMembershipsRequest,
  GetUsersForGroupInvitePickerRequest,
  InviteUserRequest,
  MembershipStatus
} from '../../api/models';
import {UserPickerBackend} from '../../shared/user-picker/user-picker.component';
import {BehaviorSubject, combineLatest, Subject} from 'rxjs';
import {GroupService} from '../group.service';

@Component({
  selector: 'app-group-members-view',
  templateUrl: './group-members-view.component.html',
  styleUrls: ['./group-members-view.component.css']
})
export class GroupMembersViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService, private groupService: GroupService) {
    this.fetchUsers = this.fetchUsers.bind(this);
  }

  refreshSubject = new Subject<boolean>();
  refresh$ = this.refreshSubject.asObservable().pipe(
    delay(500),
    startWith(true)
  );
  groupIdSubject = new BehaviorSubject(null as string);
  groupId$ = this.groupIdSubject.asObservable();
  members$ = combineLatest([this.refresh$, this.groupId$]).pipe(
    map(([_, groupId]) => groupId),
    switchMap(id => this.backend.getGroupMemberships(new GetGroupMembershipsRequest(id)))
  );
  myMembership$ = this.groupService.getMyMembership().pipe(tap((m) => console.log(m)));

  sub = this.route.parent.params.pipe(pluck('id')).subscribe(id => {
    this.groupIdSubject.next(id);
  });
  inviteUserId: string;
  pending = false;
  error = undefined;

  refreshPickerSubject = new Subject<boolean>();

  fetchUsers(skip: number, take: number, query: string) {
    return this.refreshPickerSubject.asObservable().pipe(
      delay(500),
      startWith(true)
    )
      .pipe(
        switchMap(() => {
          const getUsersQuery = new GetUsersForGroupInvitePickerRequest(skip, take, query, this.groupIdSubject.value);
          return this.backend.getUsersForGroupInvitePicker(getUsersQuery);
        }),
        pluck('users')
      );
  }

  ngOnInit(): void {

  }

  refresh() {
    this.refreshSubject.next(true);
    this.refreshPickerSubject.next(true);
  }

  inviteUser() {
    this.pending = true;
    this.error = undefined;
    this.backend.inviteUser(new InviteUserRequest(this.inviteUserId, this.groupIdSubject.value)).subscribe((res) => {
      this.pending = false;
      this.inviteUserId = null;
      this.refreshSubject.next(true);
      this.refreshPickerSubject.next(true);
    }, err => {
      this.pending = false;
      this.error = err;
    });
  }
}
