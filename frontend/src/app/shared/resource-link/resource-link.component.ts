import {Component, Input} from '@angular/core';
import {Resource, CallType} from '../../api/models';

@Component({
  selector: 'app-resource-link',
  styleUrls: ['./resource-link.component.css'],
  template: `<a [routerLink]="'/resources/' + resourceId">{{name}}</a>`
})
export class ResourceLinkComponent {

  private accountId: string;
  private callType: CallType;
  private _groupId: string;
  public resourceId: string;
  public link: string;
  public name: string;

  @Input()
  set resource(value: Resource) {
    this.callType = value?.info?.callType;
    this.accountId = value.createdBy;
    this.resourceId = value.resourceId;
    this.name = value?.info?.name;
    this.refreshLink();
  }

  @Input()
  set groupId(value: string) {
    this._groupId = value;
    this.refreshLink();
  }

  private refreshLink() {
    this.link = `/${this._groupId !== undefined ? 'groups' : 'users'}/${this._groupId !== undefined ? this._groupId : this.accountId}/${this.callType === CallType.Offer ? 'offers' : 'needs'}/${this.resourceId}`;
  }

}
