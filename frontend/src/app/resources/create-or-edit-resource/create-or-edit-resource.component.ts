import {Component, OnDestroy, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {
  CreateResourceRequest,
  GetMyMembershipsRequest,
  GetMyMembershipsResponse,
  Membership, ResourceType,
  CallType, SharedWithInput,
  UpdateResourceRequest
} from '../../api/models';
import {ActivatedRoute} from '@angular/router';
import {filter, map, pluck, shareReplay, switchMap, tap, withLatestFrom} from 'rxjs/operators';
import {AuthService} from '../../auth.service';
import {FormControl, FormGroup, Validators, FormArray} from '@angular/forms';

@Component({
  selector: 'app-new-resource',
  templateUrl: './create-or-edit-resource.component.html',
  styleUrls: ['./create-or-edit-resource.component.css']
})
export class CreateOrEditResourceComponent implements OnInit, OnDestroy {

  public submitted = false;
  public info = new FormGroup({
    callType: new FormControl(CallType.Offer, [
      Validators.required,
    ]),
    resourceType: new FormControl(ResourceType.Object, [
      Validators.required,
    ]),
    name: new FormControl('', [
      Validators.required
    ]),
    description: new FormControl('', [
      Validators.required
    ]),
  });

  public resource = new FormGroup({
    info: this.info,
    sharedWith: new FormArray([]),
    values: new FormControl(
      [{
        dimensionName: 'time',
        valueRange: {
          from: 0.3,
          to: 0.3,
        }
      }]
    ),
  });

  public form = new FormGroup({
    id: new FormControl(undefined),
    resource: this.resource
  });

  public sharedWith = this.resource.controls.sharedWith as FormArray;

  memberships$ = this.auth.session$.pipe(
    filter(s => !!s),
    pluck('id'),
    switchMap(id => this.api.getMyMemberships(new GetMyMembershipsRequest())),
    pluck<GetMyMembershipsResponse, Membership[]>('memberships'),
    map<Membership[], Membership[]>(ms => ms.filter(m => m.userConfirmed && m.groupConfirmed)),
    shareReplay()
  );

  isNewResource$ = this.route.params.pipe(pluck('id'), map(id => !id));

  resourceSub = this.route.params.pipe(pluck('id')).pipe(
    filter(id => !!id),
    switchMap(id => this.api.getResource(id)),
    shareReplay(),
  ).subscribe((res) => {
    this.sharedWith = new FormArray(res.resource.sharings.map(m => {
      return new FormGroup({
        groupId: new FormControl(m.groupId),
      });
    }));
    this.resource.setControl('sharedWith', this.sharedWith);
    const value = {
      id: res.resource.resourceId,
      resource: {
        info: {
          callType: res.resource.info.callType,
          resourceType: res.resource.info.resourceType,
          name: res.resource.info.name,
          description: res.resource.info.description,
        },
        sharedWith: res.resource.sharings.map(s => ({groupId: s.groupId})),
        values: res.resource.values
      }
    };
    this.form.setValue(value);
    this.info.controls.resourceType.disable();
    this.info.controls.callType.disable();
  });

  public error: any;
  public success = false;
  public pending = false;

  constructor(private api: BackendService, private route: ActivatedRoute, private auth: AuthService) {
  }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
    this.resourceSub.unsubscribe();
  }

  submit() {

    this.submitted = true;
    this.error = undefined;
    this.success = undefined;
    this.pending = true;

    if (this.form.value.id === null) {

      const request = CreateResourceRequest.from(this.form.value);
      this.api.createResource(request).subscribe(res => {
        this.success = true;
        this.auth.goToMyResource(res.resource.resourceId, res.resource.info.callType);
      }, err => {
        this.error = err;
        this.success = false;
        this.pending = false;
      }, () => {
        this.pending = false;
      });
    } else {
      const request = UpdateResourceRequest.from(this.form.value);
      this.api.updateResource(request).subscribe(res => {
        this.success = true;
        this.auth.goToMyResource(res.resource.resourceId, res.resource.info.callType);
      }, err => {
        this.error = err;
        this.success = false;
        this.pending = false;
      }, () => {
        this.pending = false;
      });
    }
  }

  isToggled(groupId: string) {
    for (const control of this.sharedWith.controls) {
      const isGroup = (control as FormGroup).controls.groupId.value === groupId;
      if (isGroup) {
        return true;
      }
    }
  }

  toggleGroup(groupId: string) {
    for (let i = 0; i < this.sharedWith.controls.length; i++) {
      const grpCtrl = this.sharedWith.controls[i] as FormGroup;
      const isGroup = grpCtrl.controls.groupId.value === groupId;
      if (isGroup) {
        this.sharedWith.removeAt(i);
        return;
      }
    }
    this.sharedWith.push(new FormGroup({
      groupId: new FormControl(groupId)
    }));
  }

}
