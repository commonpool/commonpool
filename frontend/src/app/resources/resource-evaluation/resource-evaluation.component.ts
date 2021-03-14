import {Component, OnDestroy, OnInit, EventEmitter} from '@angular/core';
import {ResourceEvaluationService, ResourceEvaluator} from './resource-evaluation.service';
import {empty, never, Observable, of, ReplaySubject, Subject, throwError} from 'rxjs';
import {
  finalize, last,
  map,
  mergeMap,
  retryWhen, shareReplay,
  switchMap,
  take,
  tap,
} from 'rxjs/operators';
import {BackendService} from '../../api/backend.service';
import {UpdateResourceEvaluationRequest} from '../../api/models';
import {FormControl, FormGroup} from '@angular/forms';

@Component({
  template: `
    <ng-container *ngIf="show$ | async">
      <div class="modal-backdrop fade show"></div>
      <div (click)="doNothing($event)" class="modal modal-open fade show modal-close" tabindex="-1"
           style="display:block">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">{{resourceName$ | async}}</h5>
              <button type="button" class="btn-close modal-close"></button>
            </div>
            <div class="modal-body" [formGroup]="form">
              <ng-container *ngIf="error$ | async; let error">
              <span class="text-danger">
                {{error?.message}}
              </span>
              </ng-container>
              <input type="hidden" formControlName="resourceId">
              <app-values formControlName="values"></app-values>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary modal-close">Close</button>
              <button type="button" class="btn btn-primary modal-submit">Save changes</button>
            </div>
          </div>
        </div>
      </div>
    </ng-container>
  `,
})
export class ResourceEvaluationComponent implements OnInit, OnDestroy, ResourceEvaluator {

  private showSubject = new ReplaySubject<boolean>();
  public show$ = this.showSubject.asObservable();

  private submitted = new Subject<UpdateResourceEvaluationRequest>();
  private submitted$ = this.submitted.asObservable();

  private resourceIdSubject = new ReplaySubject<string>();
  resourceId$ = this.resourceIdSubject.asObservable();

  form = new FormGroup({
    values: new FormControl(),
    resourceId: new FormControl()
  });

  public resource$ = this.resourceId$.pipe(
    switchMap(r => this.backend.getResource(r)),
    tap((r) => {
      this.form.setValue({
        resourceId: r.resource.resourceId,
        values: r.resource.values,
      });
    }),
    shareReplay()
  );
  resourceName$ = this.resource$.pipe(
    map(r => r.resource.info.name)
  );

  errorSubject = new ReplaySubject<any>();
  error$ = this.errorSubject.asObservable();

  show = false;

  public constructor(private svc: ResourceEvaluationService, private backend: BackendService) {
  }

  triggerEvaluation() {
    this.svc.evaluateResources(of(['abc', 'def'])).subscribe(value => {
      console.log('Evaluation result: ', value);
    });
  }

  doNothing($event: MouseEvent) {
    if (($event.target as HTMLElement).classList.contains('modal-close')) {
      this.submitted.next(undefined);
    } else if (($event.target as HTMLElement).classList.contains('modal-submit')) {
      this.submitted.next(UpdateResourceEvaluationRequest.from(this.form.value as UpdateResourceEvaluationRequest));
    }
    $event.preventDefault();
    $event.stopPropagation();
  }

  ngOnInit(): void {
    this.svc.setEvaluator(this);
  }

  ngOnDestroy(): void {

  }

  evaluateResource(resourceId$: Observable<string>): Observable<boolean> {
    return resourceId$.pipe(
      last(),
      switchMap((rid) => {
        return of(rid).pipe(
          tap(() => this.showSubject.next(true)),
          tap((r) => this.resourceIdSubject.next(r)),
          mergeMap(resourceId => {
            return this.submitted$.pipe(
              take(1),
              mergeMap(submitted => {
                return submitted ? of(submitted) : (throwError('aborted') as Observable<UpdateResourceEvaluationRequest>);
              })
            );
          }),
          switchMap(r => {
            return this.backend.evaluateResource$(of(r));
          }),
          map(() => true)
        );
      }),
      retryWhen(errors => errors.pipe(
        switchMap((err) => {
          const doRetry = err !== 'aborted';
          if (doRetry) {
            this.errorSubject.next(err);
            return of(true);
          }
          throw err;
        }),
      )),
      finalize(() => {
        this.showSubject.next(false);
        this.errorSubject.next(undefined);
      }),
    );
  }

}
