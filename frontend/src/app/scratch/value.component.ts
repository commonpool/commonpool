import {Component, forwardRef, Input, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {DimensionValue, ValueDimension, ValueRange} from '../api/models';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';
import {BehaviorSubject, Observable, ReplaySubject, Subject, Subscription} from 'rxjs';
import {distinctUntilChanged, map, shareReplay, tap} from 'rxjs/operators';

@Component({
  selector: 'app-value',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => ValueComponent),
      multi: true
    }
  ],
  template: `
    <input #input [disabled]="disabled"
           class="form-range"
           type="range"
           step="0.001"
           [min]="min"
           [max]="max"
           [value]="0"
           [ngModel]="(numericValue$ | async)"
           (ngModelChange)="setNumericValue($event)"
    >
  `,
})
export class ValueComponent implements ControlValueAccessor, OnDestroy, OnInit {

  private valueSubject = new ReplaySubject<ValueRange>();
  private value$: Observable<ValueRange> = this.valueSubject.asObservable()
    .pipe(distinctUntilChanged(ValueRange.equals), shareReplay());

  private valueSubscription: Subscription;
  private onChange: any;
  numericValue$ = this.value$.pipe(map(v => v?.from));
  disabled = false;

  @Input()
  min = 0;

  @Input()
  max = 0;

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
  }

  setDisabledState(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }

  setNumericValue(value: number) {
    const newValue = new ValueRange(value, value);
    this.valueSubject.next(newValue);
  }

  writeValue(obj: any): void {
    this.valueSubject.next(obj ? ValueRange.from(obj) : null);
  }

  ngOnDestroy(): void {
    this.valueSubscription.unsubscribe();
  }

  ngOnInit(): void {
    this.valueSubscription = this.value$.subscribe(value => {
      if (this.onChange) {
        this.onChange(value);
      }
    });
  }

}
