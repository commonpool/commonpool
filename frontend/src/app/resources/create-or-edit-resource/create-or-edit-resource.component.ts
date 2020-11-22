import {Component, OnDestroy, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {
  CreateResourcePayload,
  CreateResourceRequest, ExtendedResource, GetMyMembershipsRequest, GetMyMembershipsResponse,
  GetResourceResponse, Membership,
  ResourceType, SharedWithInput,
  UpdateResourcePayload,
  UpdateResourceRequest
} from '../../api/models';
import {ActivatedRoute, Router} from '@angular/router';
import {filter, map, pluck, shareReplay, switchMap, tap} from 'rxjs/operators';
import {AuthService} from '../../auth.service';
import {combineLatest} from 'rxjs';

@Component({
  selector: 'app-new-resource',
  templateUrl: './create-or-edit-resource.component.html',
  styleUrls: ['./create-or-edit-resource.component.css']
})
export class CreateOrEditResourceComponent implements OnInit, OnDestroy {

  resourceId$ = this.route.params.pipe(pluck('id'));
  resource$ = this.resourceId$.pipe(
    filter(id => !!id),
    switchMap(id => this.api.getResource(id)),
    pluck<GetResourceResponse, ExtendedResource>('resource'),
    shareReplay()
  );

  isOwnerSub = combineLatest([this.resource$, this.auth.session$]).pipe(
    map(([resource, session]) => {
      return session && resource.createdById === session.id;
    })).subscribe(isResourceOwner => {
    this.isResourceOwner = isResourceOwner;
  });

  resourceSub = this.resource$.subscribe(res => {
    this.id = res.id;
    this.summary = res.summary;
    this.description = res.description;
    this.resourceType = res.type;
    this.valueInHoursFrom = res.valueInHoursFrom;
    this.valueInHoursTo = res.valueInHoursTo;
    this.sharedWith = res.sharedWith.map(s => s.groupId);
  });

  groups$ = this.auth.session$.pipe(
    filter(s => !!s),
    pluck('id'),
    switchMap(id => this.api.getMyMemberships(new GetMyMembershipsRequest())),
    pluck<GetMyMembershipsResponse, Membership[]>('memberships'),
    map<Membership[], Membership[]>(ms => ms.filter(m => m.userConfirmed && m.groupConfirmed)),
  );

  public id: string;
  public summary: string;
  public description: string;
  public resourceType: ResourceType = ResourceType.Offer;
  public valueInHoursFrom = 1;
  public valueInHoursTo = 3;
  public isResourceOwner = true;
  public error: any;
  public success = false;
  public pending = false;
  public sharedWith: string[] = [];

  constructor(private api: BackendService, private route: ActivatedRoute, private auth: AuthService) {
  }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
    this.resourceSub.unsubscribe();
    this.isOwnerSub.unsubscribe();
  }

  submit() {

    this.error = undefined;
    this.success = undefined;
    this.pending = true;

    if (this.id === undefined) {
      const request = new CreateResourceRequest(
        new CreateResourcePayload(
          this.summary,
          this.description,
          this.resourceType,
          this.valueInHoursFrom,
          this.valueInHoursTo,
          this.sharedWith.map(s => new SharedWithInput(s))
        ));

      this.api.createResource(request).subscribe(res => {
        this.success = true;
        this.auth.goToMyResource(res.resource.id, res.resource.type);
      }, err => {
        this.error = err;
        this.success = false;
        this.pending = false;
      });
    } else {
      const request = new UpdateResourceRequest(
        this.id,
        new UpdateResourcePayload(
          this.summary,
          this.description,
          this.resourceType,
          this.valueInHoursFrom,
          this.valueInHoursTo,
          this.sharedWith.map(s => new SharedWithInput(s))
        ));

      this.api.updateResource(request).subscribe(res => {
        this.success = true;
        this.auth.goToMyResource(res.resource.id, res.resource.type);
      }, err => {
        this.error = err;
        this.success = false;
        this.pending = false;
      });
    }

  }

  toggleGroup(groupId: string) {
    console.log(groupId)
    if (this.sharedWith.includes(groupId)) {
      this.sharedWith.splice(this.sharedWith.indexOf(groupId), 1);
      this.sharedWith = [...this.sharedWith];
    } else {
      this.sharedWith.push(groupId);
      this.sharedWith = [...this.sharedWith];
    }
    console.log(this.sharedWith)
  }
}
