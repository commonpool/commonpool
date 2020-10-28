import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import {GetMembershipRequest, Membership} from '../api/models';
import {BackendService} from '../api/backend.service';
import {pluck, shareReplay} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class GroupService {

  private permissionMap: { [userKey: string]: { [groupId: string]: Observable<Membership> } } = {};

  constructor(private backend: BackendService) {
  }

  public getPermission(userId: string, groupId: string): Observable<Membership> {
    this.refreshPermission(userId, groupId);
    return this.permissionMap[userId][groupId];
  }

  public refreshPermission(userId: string, groupId: string) {
    if (!this.permissionMap[userId]) {
      this.permissionMap[userId] = {};
    }
    if (!this.permissionMap[userId][groupId]) {
      this.permissionMap[userId][groupId] = this.backend
        .getMembership(new GetMembershipRequest(userId, groupId))
        .pipe(
          pluck('membership'),
          shareReplay()
        );
    }
  }

}
