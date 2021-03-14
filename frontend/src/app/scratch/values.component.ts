import {Component, forwardRef, OnDestroy, OnInit} from '@angular/core';
import {ControlValueAccessor, FormBuilder, NG_VALUE_ACCESSOR} from '@angular/forms';
import {ValueDimensionService} from './dimension.service';
import {BehaviorSubject, combineLatest, Observable, ReplaySubject, Subject} from 'rxjs';
import {map, startWith, switchMap} from 'rxjs/operators';
import {DimensionValue, ValueDimension, ValueRange} from '../api/models';

@Component({
  selector: 'app-values',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => ValuesComponent),
      multi: true
    }
  ],
  template: `
    <div class="list-group">
      <div class="list-group-item" *ngFor="let dimension of activeDimensions$ | async">
        <div class="d-flex align-items-center">
          <div class="flex-grow-1">
            <div>{{dimension.summary}}</div>
            <div>
              <app-dimension-value
                [disabled]="disabled"
                [ngModel]="(dimensionValues$ | async)[dimension.name]"
                (ngModelChange)="valueUpdated(dimension, $event)"
              ></app-dimension-value>
            </div>
          </div>
          <div *ngIf="!disabled" class="ms-1" style="margin-right: -0.5rem">
            <button  type="button" class="btn-close" aria-label="Close"
                    (click)="removeDimension(dimension.name)"></button>
          </div>
        </div>
      </div>
      <ng-container *ngIf="!disabled && (availableDimensions$ | async); let availableDimensions">
        <div class="list-group-item" *ngIf="availableDimensions.length > 0">
          <select [disabled]="disabled" class="form-select" (change)="addDimension($event.target.value)">
            <option [value]="" selected></option>
            <option *ngFor="let dimension of availableDimensions" [value]="dimension.name">
              {{dimension.summary}}
            </option>
          </select>
        </div>
      </ng-container>
    </div>
  `,
})
export class ValuesComponent implements OnInit, OnDestroy, ControlValueAccessor {

  constructor(private svc: ValueDimensionService) {
  }

  private dimensionValuesSubject = new BehaviorSubject<{ [name: string]: DimensionValue }>({});
  public dimensionValues$: Observable<{ [name: string]: DimensionValue }> = this.dimensionValuesSubject.asObservable();
  private addDimensionSubject = new Subject<string>();
  private removeDimensionSubject = new Subject<string>();
  private onChange: any;
  disabled = false;

  addDimensionSubscription = this.addDimensionSubject.asObservable().pipe(this.svc.getDimension).subscribe(d => {
    let currentValue = this.dimensionValuesSubject.getValue();
    currentValue = {
      ...currentValue,
      ...{[d.name]: new DimensionValue(d.name, new ValueRange(d.defaultValue, d.defaultValue))}
    };
    this.dimensionValuesSubject.next(currentValue);
  });

  removeDimensionSubscription = this.removeDimensionSubject.asObservable().pipe(this.svc.getDimension).subscribe(d => {
    const currentValue = this.dimensionValuesSubject.getValue();
    const {[d.name]: _, ...newValues} = currentValue;
    this.dimensionValuesSubject.next(newValues);
  });

  private activeDimensionNames$ = this.dimensionValues$.pipe(map(v => Object.keys(v)));
  public activeDimensions$ = this.activeDimensionNames$.pipe(this.svc.getDimensionsByNames);
  public availableDimensions$ = combineLatest([this.svc.getDimensions(), this.activeDimensionNames$])
    .pipe(map(([dimensions, activeNames]) => {
      console.log(dimensions, activeNames);
      return dimensions.filter(dimension => !activeNames.includes(dimension.name));
    }));

  private fb = new FormBuilder();
  form = this.fb.array([]);

  sub = this.dimensionValues$.subscribe(v => {
    if (this.onChange) {
      const result = [];
      for (const vKey in v) {
        if (v.hasOwnProperty(vKey)) {
          result.push(v[vKey]);
        }
      }
      this.onChange(result);
    }
  });

  addDimension(name: string) {
    this.addDimensionSubject.next(name);
  }

  removeDimension(name: string) {
    this.removeDimensionSubject.next(name);
  }

  valueUpdated(dimension: ValueDimension, value: DimensionValue) {
    let currentValue = this.dimensionValuesSubject.getValue();
    currentValue = {...currentValue, ...{[dimension.name]: value}};
    this.dimensionValuesSubject.next(currentValue);
  }

  ngOnInit(): void {

  }

  ngOnDestroy(): void {
    this.addDimensionSubscription.unsubscribe();
    this.removeDimensionSubscription.unsubscribe();
    this.sub.unsubscribe();
  }

  writeValue(obj: any): void {
    const arg: DimensionValue[] = obj;
    const res = {};
    if (obj) {
      for (const dimensionValue of arg) {
        res[dimensionValue.dimensionName] = dimensionValue;
      }
    }
    this.dimensionValuesSubject.next(res);
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {

  }

  setDisabledState?(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }

}
