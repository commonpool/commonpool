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
    value: new FormGroup({
      timeValueFrom: new FormControl(0, [
        Validators.required,
        Validators.min(0)
      ]),
      timeValueTo: new FormControl(0, [
        Validators.required,
        Validators.min(0),
      ]),
    })

  });

  public resource = new FormGroup({
    info: this.info,
    sharedWith: new FormArray([])
  }, control => {
    const fg = control as FormGroup;
    if (fg.controls.timeValueFrom > fg.controls.timeValueTo) {
      return {hoursFromLargerThanValuesTo: {}};
    }
  });

  public form = new FormGroup({
    id: new FormControl(undefined),
    resource: this.resource
  });

  public sharedWith = this.resource.controls.sharedWith as FormArray;

  formValueChanged = this.form.valueChanges.subscribe((v) => {
    const resourceControls = (this.form.controls.resource as FormGroup);
    if (v.resource.timeValueFrom < 0) {
      resourceControls.controls.timeValueFrom.setValue(0);
    } else if (v.resource.timeValueTo < 0) {
      resourceControls.controls.timeValueTo.setValue(0);
    } else if (v.resource.timeValueFrom > v.resource.timeValueTo) {
      resourceControls.controls.timeValueFrom.setValue(v.resource.timeValueTo);
    }
  });

  memberships$ = this.auth.session$.pipe(
    filter(s => !!s),
    pluck('id'),
    switchMap(id => this.api.getMyMemberships(new GetMyMembershipsRequest())),
    pluck<GetMyMembershipsResponse, Membership[]>('memberships'),
    map<Membership[], Membership[]>(ms => ms.filter(m => m.userConfirmed && m.groupConfirmed)),
    shareReplay()
  );

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
          value: {
            timeValueFrom: res.resource.info.value.timeValueFrom / 1000000000,
            timeValueTo: res.resource.info.value.timeValueTo / 1000000000,
          }
        },
        sharedWith: res.resource.sharings.map(s => ({groupId: s.groupId}))
      }
    };
    this.form.setValue(value);
    this.resource.controls.type.disable();
    this.resource.controls.subType.disable();
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
    this.formValueChanged.unsubscribe();
  }

  submit() {

    this.submitted = true;
    this.error = undefined;
    this.success = undefined;
    this.pending = true;

    const value = this.form.value;
    value.resource.info.value.timeValueFrom *= 1000000000;
    value.resource.info.value.timeValueTo *= 1000000000;

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
