import {
  AbstractControl,
  AbstractControlOptions, AsyncValidatorFn,
  FormArray,
  FormControl,
  FormGroup,
  ValidatorFn,
  Validators
} from '@angular/forms';
import {OfferItemType, Target} from '../../api/models';
import {distinctUntilChanged, pluck} from 'rxjs/operators';

export const minLengthArray = (min: number) => {
  return (c: AbstractControl): { [key: string]: any } => {
    if (c.value.length >= min) {
      return null;
    }
    return {MinLengthArray: `must have at least ${min} items`};
  };
};

export const uuidValidator = (property?: string) => {
  return (c: AbstractControl): { [key: string]: any } => {
    let value = c.value;
    if (property) {
      value = value[property];
    }
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
  public groupId: FormControl;

  public constructor() {
    super({
      groupId: new FormControl('', [Validators.required]),
      items: new FormArray([], [minLengthArray(1)]),
      message: new FormControl('')
    });
    this.items = this.controls.items as FormArray;
    this.message = this.controls.message as FormControl;
    this.groupId = this.controls.groupId as FormControl;
  }

  public removeItem(i: number) {
    this.items.removeAt(i);
  }

  public getItem(i: number): CreateOfferItemForm {
    return this.items.controls[i] as CreateOfferItemForm;
  }

  public getItems(): CreateOfferItemForm[] {
    return this.items.controls as CreateOfferItemForm[];
  }

  getErrors(): any {
    return getFormErrors(this);
  }
}

export class TargetForm extends FormControl {

  constructor(
    target: Target,
    validatorOrOpts?: ValidatorFn | ValidatorFn[] | AbstractControlOptions | null,
    asyncValidator?: AsyncValidatorFn | AsyncValidatorFn[] | null) {
    super(undefined, validatorOrOpts, asyncValidator);
  }

}

export class CreateOfferItemForm extends FormGroup {

  public fromControl = new TargetForm(undefined, []);
  public toControl = new TargetForm(undefined, [Validators.required, (c) => {

    if (c.value && c.value.type === 'group') {
      return uuidValidator('groupId');
    } else if (c.value && c.value.type === 'user') {
      return uuidValidator('userId');
    }

  }]);
  public typeControl = new FormControl(undefined, [Validators.required]);
  public resourceIdControl = new FormControl('');
  public amountControl = new FormControl('');
  public durationControl = new FormControl('');

  private readonly fromKey = 'from';
  private readonly toKey = 'to';
  private readonly typeKey = 'type';
  private readonly resourceIdKey = 'resourceId';
  private readonly amountKey = 'amount';
  private readonly durationKey = 'duration';

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
    const type = this.getType();
    if (type === OfferItemType.BorrowResource
      || type === OfferItemType.ProvideService
      || type === OfferItemType.ResourceTransfer) {
      this.amountControl.setValidators([]);
      this.resourceIdControl.setValidators([
        Validators.required
      ]);
    } else {
      this.resourceIdControl.setValidators([]);
      this.amountControl.setValidators([
        Validators.required,
        Validators.min(0)
      ]);
    }
    this.amountControl.updateValueAndValidity();
    this.resourceIdControl.updateValueAndValidity();
  }

  constructor() {
    super({});
    this.addControl(this.fromKey, this.fromControl);
    this.addControl(this.toKey, this.toControl);
    this.addControl(this.typeKey, this.typeControl);
    this.addControl(this.resourceIdKey, this.resourceIdControl);
    this.addControl(this.amountKey, this.amountControl);
    this.addControl(this.durationKey, this.durationControl);
    this.updateValidators();
  }

  getErrors(): any {
    return {
      from: {...this.fromControl.errors},
      to: {...this.toControl.errors},
      type: {...this.typeControl.errors},
      resourceId: {...this.resourceIdControl.errors},
      amount: {...this.amountControl.errors},
      duration: {...this.durationControl.errors},
    };
  }
}
