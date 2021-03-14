import {Injectable} from '@angular/core';
import {Observable, of, ReplaySubject, Subject} from 'rxjs';
import {catchError, concatMap, map, reduce, scan, switchMap, take, tap, withLatestFrom} from 'rxjs/operators';
import {fromArray} from 'rxjs/internal/observable/fromArray';

export interface ResourceEvaluator {
  evaluateResource(resourceId$: Observable<string>): Observable<boolean>;
}

@Injectable({
  providedIn: 'root'
})
export class ResourceEvaluationService {

  private evaluatorSubject = new ReplaySubject<ResourceEvaluator>();
  private evaluator$ = this.evaluatorSubject.asObservable();

  public constructor() {
    this.evaluateResource = this.evaluateResource.bind(this);
    this.evaluateResources = this.evaluateResources.bind(this);
  }

  public setEvaluator(evaluator: ResourceEvaluator) {
    this.evaluatorSubject.next(evaluator);
  }

  evaluateResource(resourceId$: Observable<string>): Observable<boolean> {
    return resourceId$.pipe(
      take(1),
      tap((rid) => console.log('resource id', rid)),
      withLatestFrom(this.evaluator$),
      switchMap(([resourceId, evaluator]) => {
        return evaluator.evaluateResource(of(resourceId));
      }),
      catchError(() => of(false))
    );
  }

  evaluateResources(resourceIds$: Observable<string[]>): Observable<boolean> {
    const result = resourceIds$.pipe(
      switchMap(resourceIds => {
        return resourceIds.reduce((previous, next) => {
          return previous.pipe(switchMap(v => {
            if (!v) {
              return of(false);
            }
            return this.evaluateResource(of(next));
          }));
        }, of(true));
      }),
      tap(console.log),
      )
    ;
    return result;
  }

}
