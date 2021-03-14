import {Component, EventEmitter, forwardRef, Input, OnDestroy, OnInit} from '@angular/core';
import {BackendService} from '../api/backend.service';
import {DimensionValue, GetValueDimensionsRequest, ValueDimension, ValueRange} from '../api/models';
import {combineLatest, Observable, ReplaySubject, Subscription} from 'rxjs';
import {distinctUntilChanged, filter, map, shareReplay} from 'rxjs/operators';
import {ControlValueAccessor, FormArray, FormControl, NG_VALUE_ACCESSOR} from '@angular/forms';
import {ValueDimensionService} from './dimension.service';

@Component({
  selector: 'app-dimension-value',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => DimensionValueComponent),
      multi: true
    }
  ],
  template: `
    <ng-container *ngIf="dimension$ | async">
      <app-value
        [min]="from$ | async"
        [max]="to$ | async"
        [disabled]="disabled"
        [ngModel]="valueRange$ | async"
        (ngModelChange)="valueRangeSubject.next($event)"
      ></app-value>
      <app-value-threshold
        [min]="from$ | async"
        [max]="to$ | async"
        [thresholds]="thresholds$ | async"
        [value]="valueRange$ | async"></app-value-threshold>
    </ng-container>
  `,
})
export class DimensionValueComponent implements OnInit, OnDestroy, ControlValueAccessor {

  private _onChange: any;

  public constructor(private dimensionService: ValueDimensionService) {

  }

  valueRangeSubject = new ReplaySubject<ValueRange>();
  valueRange$: Observable<ValueRange> = this.valueRangeSubject.asObservable();
  private dimensionNameSubject = new ReplaySubject<string>();
  private dimensionName$ = this.dimensionNameSubject.asObservable();
  dimension$ = this.dimensionName$.pipe(this.dimensionService.getDimension);
  from$ = this.dimension$.pipe(map(d => d.range.from));
  to$ = this.dimension$.pipe(map(d => d.range.to));
  thresholds$ = this.dimension$.pipe(map(d => d.thresholds));
  private value$ = combineLatest([this.dimensionName$, this.valueRange$]).pipe(
    map(([dimensionName, valueRange]) => new DimensionValue(dimensionName, valueRange)),
    distinctUntilChanged(DimensionValue.equals),
    shareReplay()
  );
  private valueSubscription: Subscription;
  public formReady: EventEmitter<null> = new EventEmitter<null>();
  disabled = false

  ngOnInit(): void {
    this.valueSubscription = this.value$.subscribe((v) => {
      if (this._onChange) {
        this._onChange(v);
      }
    });
  }

  ngOnDestroy(): void {
    this.valueSubscription.unsubscribe();
  }

  writeValue(obj: any): void {
    if (obj) {
      const value = DimensionValue.from(obj);
      this.dimensionNameSubject.next(value.dimensionName);
      this.valueRangeSubject.next(value.valueRange);
    }
  }

  registerOnChange(fn: any): void {
    this._onChange = fn;
  }

  registerOnTouched(fn: any): void {

  }

  setDisabledState?(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }


}
