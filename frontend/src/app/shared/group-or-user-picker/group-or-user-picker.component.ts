import {Component, forwardRef, Input, OnInit} from '@angular/core';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';
import {combineLatest, of, ReplaySubject, Subject} from 'rxjs';
import {OfferItemTargetRequest, OfferItemType, Target} from '../../api/models';
import {BackendService} from '../../api/backend.service';
import {
  debounceTime,
  distinctUntilChanged,
  filter, map,
  shareReplay,
  startWith,
  switchMap,
  tap,
  withLatestFrom
} from 'rxjs/operators';

@Component({
  selector: 'app-group-or-user-picker',
  templateUrl: './group-or-user-picker.component.html',
  styleUrls: ['./group-or-user-picker.component.css'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => GroupOrUserPickerComponent),
      multi: true
    }
  ]
})
export class GroupOrUserPickerComponent implements OnInit, ControlValueAccessor {

  private propagateChangeFn: (val: Target) => void;

  private groupIdSubject = new ReplaySubject<string>();
  private offerItemTypeSubject = new ReplaySubject<OfferItemType>();
  public selectedTargetSubject = new ReplaySubject<Target>();
  public selectedTarget$ = this.selectedTargetSubject.asObservable();
  private toTarget = new ReplaySubject<Target>();
  selectedTargetSub = this.selectedTarget$.subscribe(target => {
    if (this.propagateChangeFn) {
      this.propagateChangeFn(target);
    }
  });

  bla = Math.random();

  public targets$ = combineLatest([
    this.groupIdSubject.asObservable().pipe(distinctUntilChanged()),
    this.offerItemTypeSubject.pipe(distinctUntilChanged()),
    this.toTarget.asObservable().pipe(
      startWith(undefined as Target),
      distinctUntilChanged((x, y) => {
        return (x?.type === y?.type && x?.userId === y?.userId && x?.groupId === y?.groupId);
      })
    )])
    .pipe(
      debounceTime(50),
      filter(([groupId, offerItemType, _]) => !!groupId && !!offerItemType),
      switchMap(([groupId, offerItemType, to]) => {
          let fromId;
          if (to) {
            if (to.type === 'group') {
              fromId = to.groupId;
            } else if (to.type === 'user') {
              fromId = to.userId;
            }
          }
          return this.backend.getItemsForTargetPicker(new OfferItemTargetRequest(
            offerItemType,
            groupId,
            undefined,
            undefined,
            to ? to.type : undefined,
            fromId));
        }
      ),
      withLatestFrom(this.selectedTarget$),
      tap(([targets, selected]) => {

        if (!selected) {
          return;
        }
        const found = targets.items.find(target => {
          return (target.type === selected.type && target.groupId === selected.groupId && target.userId === selected.userId);
        });
        if (!found) {
          this.selectedTargetSubject.next(undefined);
        }
      }),
      map(([targets, _]) => targets),
      shareReplay()
    );

  @Input()
  set groupId(groupId: string) {
    this.groupIdSubject.next(groupId);
  }

  @Input()
  set offerItemType(offerItemType: OfferItemType) {
    this.offerItemTypeSubject.next(offerItemType);
  }

  @Input()
  set to(target: Target) {
    this.toTarget.next(target);
  }

  constructor(private backend: BackendService) {
  }

  ngOnInit(): void {
  }

  registerOnChange(fn: any): void {
    this.propagateChangeFn = fn;
  }

  registerOnTouched(fn: any): void {
  }

  setDisabledState(isDisabled: boolean): void {
  }

  writeValue(obj: any): void {
    this.selectedTargetSubject.next(obj);
  }

}
