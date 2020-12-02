import {Component, OnDestroy} from '@angular/core';
import {OfferItemType, SendOfferRequest} from '../../api/models';
import {BackendService} from '../../api/backend.service';
import {CreateOfferForm, CreateOfferItemForm} from './create-offer.form';
import {distinctUntilChanged, pluck} from 'rxjs/operators';

@Component({
  selector: 'app-create-offer',
  templateUrl: './create-offer.component.html',
  styleUrls: ['./create-offer.component.css']
})
export class CreateOfferComponent implements OnDestroy {

  constructor(private backend: BackendService) {
  }

  public form = new CreateOfferForm();
  public itemForm = new CreateOfferItemForm();

  submitted = false;
  pending = false;
  error = undefined;
  offerItemType: OfferItemType = OfferItemType.ResourceItem;

  formValueSub = this.itemForm.valueChanges.pipe(
    pluck<any, string>('from'),
    distinctUntilChanged(),
  ).subscribe((fromUserId: string) => {
    const predicate = (toUserId: string) => {
      console.log(toUserId, fromUserId);
      return toUserId !== fromUserId;
    };
    predicate.bind(this);
    this.toPredicate = predicate;
  });

  toPredicate = (val: string) => true;

  add() {
    const newItemForm = new CreateOfferItemForm();
    let resourceId = this.itemForm.resourceIdControl.value;
    if (!resourceId) {
      resourceId = '';
    }
    newItemForm.setValue({
      ...this.itemForm.value,
      resourceId,
      timeInSeconds: this.itemForm.value.timeInSeconds * 60 * 60
    });
    this.form.items.push(newItemForm);
    this.itemForm.setParent(newItemForm);
    this.itemForm.setValue({
      from: '',
      to: '',
      type: OfferItemType.ResourceItem,
      resourceId: '',
      timeInSeconds: 0
    });
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
      console.log(res);
    }, err => {
      this.pending = false;
      this.error = err;
    });
  }

  ngOnDestroy(): void {
    if (this.formValueSub) {
      this.formValueSub.unsubscribe();
    }
  }

}
