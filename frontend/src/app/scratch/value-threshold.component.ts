import {Component, Input, OnInit} from '@angular/core';
import {ValueDimension, ValueRange, ValueThreshold} from '../api/models';
import {combineLatest, Observable, ReplaySubject} from 'rxjs';
import {map} from 'rxjs/operators';

@Component({
  selector: 'app-value-threshold',
  template: `
    {{(threshold$ | async)}}
  `,
})
export class ValueThresholdComponent {

  private valueSubject = new ReplaySubject<ValueRange>();
  private minSubject = new ReplaySubject<number>();
  private maxSubject = new ReplaySubject<number>();
  private thresholdsSubject = new ReplaySubject<ValueThreshold[]>();
  private thresholds$ = this.thresholdsSubject.asObservable();
  min$ = this.minSubject.asObservable();
  max$ = this.maxSubject.asObservable();
  value$: Observable<ValueRange> = this.valueSubject.asObservable();

  public threshold$ = combineLatest([
    this.min$,
    this.max$,
    this.value$,
    this.thresholds$]).pipe(
    map(([min, max, value, thresholds]) => {
      const range = max - min;
      const step = range / (thresholds.length - 1);
      const abs = - min + value.from;
      const idx = abs / step;
      const snap = Math.round(idx);
      return thresholds[snap].description;
    }));

  constructor() {
  }

  @Input()
  set value(value: ValueRange) {
    this.valueSubject.next(value);
  }

  @Input()
  set thresholds(value: ValueThreshold[]) {
    this.thresholdsSubject.next(value);
  }

  @Input()
  set min(value: number) {
    this.minSubject.next(value);
  }

  @Input()
  set max(value: number) {
    this.maxSubject.next(value);
  }

}
