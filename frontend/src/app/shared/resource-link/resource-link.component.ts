import {Component, Input} from '@angular/core';
import {Resource, ResourceType} from '../../api/models';

@Component({
  selector: 'app-resource-link',
  styleUrls: ['./resource-link.component.css'],
  template: ` <a
    [routerLink]="link">{{summary}}</a>`
})
export class ResourceLinkComponent {

  private accountId: string;
  private resourceType: ResourceType;
  private accountType = 'user';
  private _groupId: string;
  private _resourceId: string;
  public link: string;
  public summary: string;

  @Input()
  set resource(value: Resource) {
    this.resourceType = value.type;
    this.accountId = value.createdById;
    this._resourceId = value.id;
    this.summary = value.summary;
    this.refreshLink();
  }

  @Input()
  set groupId(value: string) {
    this._groupId = value;
    this.refreshLink();
  }

  private refreshLink() {
    this.link = `/${this._groupId !== undefined ? 'groups' : 'users'}/${this._groupId !== undefined ? this._groupId : this.accountId}/${this.resourceType === ResourceType.Offer ? 'offers' : 'needs'}/${this._resourceId}`;
  }

}
