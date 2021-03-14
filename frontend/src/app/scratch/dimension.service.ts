import {Injectable} from '@angular/core';
import {BackendService} from '../api/backend.service';
import {Observable, ReplaySubject} from 'rxjs';
import {GetValueDimensionsRequest, ValueDimension} from '../api/models';
import {map, pluck, shareReplay, switchMap, tap} from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class ValueDimensionService {

  public dimensions$ = this.backend
    .getValueDimensions(new GetValueDimensionsRequest())
    .pipe(pluck('dimensions'), shareReplay());

  public constructor(private backend: BackendService) {
    this.getDimensions = this.getDimensions.bind(this);
    this.getDimension = this.getDimension.bind(this);
    this.getDimensionsByNames = this.getDimensionsByNames.bind(this);
  }

  public getDimensions(): Observable<ValueDimension[]> {
    return this.dimensions$;
  }

  public getDimension(name: Observable<string>): Observable<ValueDimension> {
    return name.pipe(switchMap(dimensionName => {
      return this.getDimensions().pipe(map(dimensions => {
        return dimensions.find(d => d.name === dimensionName);
      }));
    }));
  }

  public getDimensionsByNames(names: Observable<string[]>): Observable<ValueDimension[]> {
    return names.pipe(switchMap(dimensionNames => {
      return this.getDimensions().pipe(map(ds => {
        return ds.filter(d => dimensionNames.includes(d.name));
      }));
    }));
  }
}


