import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {pluck, switchMap} from 'rxjs/operators';
import {BackendService} from '../../api/backend.service';
import {GetGroupMembershipsRequest, GetUsersForGroupInvitePickerRequest, InviteUserRequest} from '../../api/models';
import {UserPickerBackend} from '../../shared/user-picker/user-picker.component';
import {BehaviorSubject} from 'rxjs';

@Component({
  selector: 'app-group-members-view',
  templateUrl: './group-members-view.component.html',
  styleUrls: ['./group-members-view.component.css']
})
export class GroupMembersViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService) {
    this.fetchUsers = this.fetchUsers.bind(this);
  }

  groupIdSubject = new BehaviorSubject(null as string);
  groupId$ = this.groupIdSubject.asObservable();
  members$ = this.groupId$.pipe(
    switchMap(id => this.backend.getGroupMemberships(new GetGroupMembershipsRequest(id)))
  );
  sub = this.route.parent.params.pipe(pluck('id')).subscribe(id => {
    this.groupIdSubject.next(id);
  });
  inviteUserId: string;
  pending = false;
  error = undefined;

  fetchUsers(skip: number, take: number, query: string) {
    return this.backend.getUsersForGroupInvitePicker(new GetUsersForGroupInvitePickerRequest(skip, take, query, this.groupIdSubject.value))
      .pipe(pluck('users'));
  }


  ngOnInit(): void {

  }

  inviteUser() {
    this.pending = true;
    this.error = undefined;
    this.backend.inviteUser(new InviteUserRequest(this.inviteUserId, this.groupIdSubject.value)).subscribe((res) => {
      this.pending = false;
      this.inviteUserId = null;
    }, err => {
      this.pending = false;
      this.error = err;
    });
  }
}
