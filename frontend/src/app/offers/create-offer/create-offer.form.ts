import {AbstractControl, FormArray, FormControl, FormGroup, ValidationErrors, Validators} from '@angular/forms';
import {Offer, OfferItemType} from '../../api/models';
import {distinctUntilChanged, pluck} from 'rxjs/operators';

export const minLengthArray = (min: number) => {
  return (c: AbstractControl): { [key: string]: any } => {
    if (c.value.length >= min) {
      return null;
    }
    return {MinLengthArray: `must have at least ${min} items`};
  };
};

export const uuidValidator = () => {
  return (c: AbstractControl): { [key: string]: any } => {
    const value = c.value;
    const regex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
    if (!regex.test(value)) {
      return {
        uuid: 'invalid uuid'
      };
    } else {
      return null;
    }
  };
};

function getFormErrors(form: AbstractControl) {
  if (form instanceof FormControl) {
    // Return FormControl errors or null
    return form.errors ?? null;
  }
  if (form instanceof FormGroup || form instanceof FormArray) {
    const groupErrors = form.errors;
    // Form group can contain errors itself, in that case add'em
    const formErrors = groupErrors ? {groupErrors} : {};
    Object.keys(form.controls).forEach(key => {
      // Recursive call of the FormGroup fields
      const error = getFormErrors(form.get(key));
      if (error !== null) {
        // Only add error if not null
        formErrors[key] = error;
      }
    });
    // Return FormGroup errors or null
    return Object.keys(formErrors).length > 0 ? formErrors : null;
  }
}

export class CreateOfferForm extends FormGroup {

  public items: FormArray;
  public message: FormControl;

  public constructor() {
    super({
      items: new FormArray([], [minLengthArray(1)]),
      message: new FormControl('')
    });
    this.items = this.controls.items as FormArray;
    this.message = this.controls.message as FormControl;
  }

  public removeItem(i: number) {
    this.items.removeAt(i);
  }

  public getItem(i: number): CreateOfferItemForm {
    return this.items.controls[i] as CreateOfferItemForm;
  }

  getErrors(): any {
    return getFormErrors(this);
  }
}

export class CreateOfferItemForm extends FormGroup {

  public fromControl = new FormControl('', [Validators.required, uuidValidator()]);
  public toControl = new FormControl('', [Validators.required, uuidValidator()]);
  public typeControl = new FormControl(OfferItemType.ResourceItem, [Validators.required]);
  public resourceIdControl = new FormControl('');
  public timeInSecondsControl = new FormControl(0);

  private readonly fromKey = 'from';
  private readonly toKey = 'to';
  private readonly typeKey = 'type';
  private readonly resourceIdKey = 'resourceId';
  private readonly timeInSecondsKey = 'timeInSeconds';

  private valueSub = this.valueChanges.pipe(
    pluck('type'),
    distinctUntilChanged()
  ).subscribe(v => {
    this.updateValidators();
  });

  public getType(): OfferItemType {
    return this.typeControl.value as OfferItemType;
  }

  private updateValidators() {
    if (this.getType() === OfferItemType.ResourceItem) {
      this.timeInSecondsControl.setValidators([]);
      this.resourceIdControl.setValidators([
        Validators.required
      ]);
    } else {
      this.resourceIdControl.setValidators([]);
      this.timeInSecondsControl.setValidators([
        Validators.required,
        Validators.min(0)
      ]);
    }
    this.timeInSecondsControl.updateValueAndValidity();
    this.resourceIdControl.updateValueAndValidity();
  }

  constructor() {
    super({});
    this.addControl(this.fromKey, this.fromControl);
    this.addControl(this.toKey, this.toControl);
    this.addControl(this.typeKey, this.typeControl);
    this.addControl(this.resourceIdKey, this.resourceIdControl);
    this.addControl(this.timeInSecondsKey, this.timeInSecondsControl);
    this.updateValidators();
  }

  getErrors(): any {
    return {
      from: {...this.fromControl.errors},
      to: {...this.toControl.errors},
      type: {...this.typeControl.errors},
      resourceId: {...this.resourceIdControl.errors},
      timeInSeconds: {...this.timeInSecondsControl.errors}
    };
  }
}
