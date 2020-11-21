import {Component, Input} from '@angular/core';
import {Resource, ResourceType} from '../../api/models';

@Component({
  selector: 'app-resource-link',
  styleUrls: ['./resource-link.component.css'],
  template: ` <a
    [routerLink]="link">{{resource.summary}}</a>`
})
export class ResourceLinkComponent {

  private _resource: Resource;
  public link: string;

  @Input()
  set resource(value: Resource) {
    this._resource = value;
    this.link = `/users/${value.createdBy}/${value.type === ResourceType.Offer ? 'offers' : 'needs'}/${value.id}`;
  }

  get resource(): Resource {
    return this._resource;
  }

}
