import {combineLatest, Observable, OperatorFunction, Subject} from 'rxjs';
import {map, startWith} from 'rxjs/operators';

export class RefreshTrigger {

  private refreshSubject = new Subject<boolean>();
  private refresh$ = this.refreshSubject.asObservable().pipe(startWith(true));

  constructor() {
  }

  trigger() {
    this.refreshSubject.next(true);
  }

  public asObservable() {
    return this.refresh$;
  }

  public refresh<T>(sub: Observable<T>) {
    return combineLatest([
      this.asObservable(),
      sub
    ]).pipe(
      map(([_, value]) => value)
    );
  }

}

export function refreshWith(refreshSubject: Subject<any>) {
  return <T>(source: Observable<T>): Observable<T> => {
    let lastValue;
    let calledOnce = false;
    return new Observable(subscriber => {
      refreshSubject.subscribe(value => {
        if (calledOnce) {
          subscriber.next(lastValue);
        }
      }, error => {
        subscriber.error(error);
      }, () => {
        subscriber.complete();
      });
      source.subscribe(value => {
        lastValue = value;
        calledOnce = true;
        subscriber.next(value);
      }, error => {
        subscriber.error(error);
      }, () => {
        subscriber.complete();
      });
    });
  };
}
