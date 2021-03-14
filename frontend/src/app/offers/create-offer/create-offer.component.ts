import {Component, OnDestroy} from '@angular/core';
import {OfferItemType, SendOfferRequest} from '../../api/models';
import {BackendService} from '../../api/backend.service';
import {CreateOfferForm, CreateOfferItemForm} from './create-offer.form';
import {distinctUntilChanged, pluck} from 'rxjs/operators';
import {Router} from '@angular/router';
import {AuthService} from '../../auth.service';

@Component({
  selector: 'app-create-offer',
  templateUrl: './create-offer.component.html',
  styleUrls: ['./create-offer.component.css']
})
export class CreateOfferComponent implements OnDestroy {
  constructor(private backend: BackendService, private router: Router, private auth: AuthService) {
  }

  public form = new CreateOfferForm();
  public itemForm = new CreateOfferItemForm();
  submitted = false;
  pending = false;
  error = undefined;
  itemFormToggled = true;
  groupSelected = false;
  disableGroupSub = this.form.valueChanges.subscribe(value => {
    if (value.type && this.form.groupId.enabled) {
      this.form.groupId.disable();
    } else if (!value.type && this.form.groupId.disabled) {
      this.form.groupId.enable();
    }
  });
  formValueSub = this.itemForm.valueChanges.pipe(
    pluck<any, string>('from'),
    distinctUntilChanged(),
  ).subscribe((fromUserId: string) => {
    const predicate = (toUserId: string) => {
      return toUserId !== fromUserId;
    };
    predicate.bind(this);
    this.toPredicate = predicate;
  });

  markGroupSelected() {
    this.groupSelected = true;
    this.form.groupId.disable();
  }

  toPredicate = (val: string) => true;

  add() {
    const newItemForm = new CreateOfferItemForm();
    let resourceId = this.itemForm.resourceIdControl.value;
    if (!resourceId) {
      resourceId = null;
    }
    newItemForm.setValue({
      ...this.itemForm.value,
      resourceId,
    });
    this.form.items.push(newItemForm);
    this.itemForm.setParent(newItemForm);
    this.itemForm.setValue({
      from: null,
      to: null,
      type: OfferItemType.BorrowResource,
      resourceId: null,
      amount: null,
      duration: null
    });
    this.itemFormToggled = false;
  }

  remove(i: number) {
    this.form.removeItem(i);
    if (this.form.items.controls.length === 0) {
      this.itemFormToggled = true;
    }
  }

  submit() {
    this.submitted = true;

    if (!this.form.valid) {
      return;
    }

    this.pending = true;
    this.error = undefined;

    const request = {offer: this.form.value} as SendOfferRequest;

    this.backend.sendOffer(SendOfferRequest.from(request)).subscribe(res => {
      this.pending = false;
      this.router.navigateByUrl('/users/' + this.auth.getCurrentAuthId() + '/transactions');
    }, err => {
      this.pending = false;
      this.error = err;
    });
  }

  ngOnDestroy(): void {
    this.formValueSub.unsubscribe();
    this.disableGroupSub.unsubscribe();
  }

}
