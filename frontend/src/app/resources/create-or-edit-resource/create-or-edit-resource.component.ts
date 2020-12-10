import {Component, OnDestroy, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {
  CreateResourceRequest,
  GetMyMembershipsRequest,
  GetMyMembershipsResponse,
  Membership, ResourceSubType,
  ResourceType, SharedWithInput,
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

  public resource = new FormGroup({
    type: new FormControl(ResourceType.Offer, [
      Validators.required,
    ]),
    subType: new FormControl(ResourceSubType.Object, [
      Validators.required,
    ]),
    summary: new FormControl('', [
      Validators.required
    ]),
    description: new FormControl('', [
      Validators.required
    ]),
    valueInHoursFrom: new FormControl(0, [
      Validators.required,
      Validators.min(0)
    ]),
    valueInHoursTo: new FormControl(0, [
      Validators.required,
      Validators.min(0),
    ]),
    sharedWith: new FormArray([])
  }, control => {
    const fg = control as FormGroup;
    if (fg.controls.valueInHoursFrom > fg.controls.valueInHoursTo) {
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
    if (v.resource.valueInHoursFrom < 0) {
      resourceControls.controls.valueInHoursFrom.setValue(0);
    } else if (v.resource.valueInHoursTo < 0) {
      resourceControls.controls.valueInHoursTo.setValue(0);
    } else if (v.resource.valueInHoursFrom > v.resource.valueInHoursTo) {
      resourceControls.controls.valueInHoursFrom.setValue(v.resource.valueInHoursTo);
    }
  });

  memberships$ = this.auth.session$.pipe(
    filter(s => !!s),
    pluck('id'),
    switchMap(id => this.api.getMyMemberships(new GetMyMembershipsRequest())),
    pluck<GetMyMembershipsResponse, Membership[]>('memberships'),
    map<Membership[], Membership[]>(ms => ms.filter(m => m.userConfirmed && m.groupConfirmed)),
    tap(memberships => {

      console.log(this.resource.controls.sharedWith);
    }),
    shareReplay()
  );

  resourceSub = this.route.params.pipe(pluck('id')).pipe(
    filter(id => !!id),
    switchMap(id => this.api.getResource(id)),
    shareReplay(),
  ).subscribe((res) => {
    this.sharedWith = new FormArray(res.resource.sharedWith.map(m => {
      return new FormGroup({
        groupId: new FormControl(m.groupId),
      });
    }));
    this.resource.setControl('sharedWith', this.sharedWith);
    const value = {
      id: res.resource.id,
      resource: {
        type: res.resource.type,
        subType: res.resource.subType,
        summary: res.resource.summary,
        description: res.resource.description,
        valueInHoursFrom: res.resource.valueInHoursFrom,
        valueInHoursTo: res.resource.valueInHoursTo,
        sharedWith: res.resource.sharedWith.map(s => ({groupId: s.groupId}))
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

    if (this.form.value.id === null) {
      const request = CreateResourceRequest.from(this.form.value);
      this.api.createResource(request).subscribe(res => {
        this.success = true;
        this.auth.goToMyResource(res.resource.id, res.resource.type);
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
        this.auth.goToMyResource(res.resource.id, res.resource.type);
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
