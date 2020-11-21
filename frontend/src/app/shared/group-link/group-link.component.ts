import {Component, Input, OnInit} from '@angular/core';
import {Group, Membership} from '../../api/models';

export class GroupOrMembership {
  public groupName: string;
  public groupId: string;

  public constructor(public groupOrMembership: Group | Membership) {
    if (groupOrMembership instanceof Group) {
      this.groupName = (groupOrMembership as Group).name;
      this.groupId = (groupOrMembership as Group).id;
    } else if (groupOrMembership instanceof Membership) {
      this.groupName = (groupOrMembership as Membership).groupName;
      this.groupId = (groupOrMembership as Membership).groupId;
    }
  }
}

@Component({
  selector: 'app-group-link',
  template: `<a [routerLink]="'/groups/' + _groupOrMembership.groupId">{{_groupOrMembership.groupName}}</a>`,
})
export class GroupLinkComponent {

  _groupOrMembership: GroupOrMembership;

  @Input()
  set groupOrMembership(groupOrMembership: Group | Membership) {
    this._groupOrMembership = new GroupOrMembership(groupOrMembership);
  }

}
