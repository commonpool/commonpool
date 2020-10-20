import {Component, OnDestroy, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {CreateResourcePayload, CreateResourceRequest, ResourceType, UpdateResourcePayload, UpdateResourceRequest} from '../../api/models';
import {ActivatedRoute, Router} from '@angular/router';
import {filter, map, pluck, shareReplay, switchMap} from 'rxjs/operators';
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
    pluck('resource'),
    shareReplay()
  );

  isOwnerSub = combineLatest([this.resource$, this.auth.session$]).pipe(map(([resource, session]) => {
    return resource.createdById === session.id;
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
  });

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
          this.valueInHoursTo
        ));

      this.api.createResource(request).subscribe(res => {
        this.success = true;
        this.auth.goToMyProfile();
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
          this.valueInHoursTo
        ));

      this.api.updateResource(request).subscribe(res => {
        this.success = true;
        this.auth.goToMyProfile();
      }, err => {
        this.error = err;
        this.success = false;
        this.pending = false;
      });
    }


  }


}
