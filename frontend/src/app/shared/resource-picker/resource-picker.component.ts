import {Component, forwardRef, Input, OnDestroy, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {combineLatest, of, ReplaySubject, Subject, Subscription} from 'rxjs';
import {distinctUntilChanged, map, pluck, shareReplay, startWith, switchMap, tap} from 'rxjs/operators';
import {Resource, ResourceType, CallType, SearchResourceRequest} from '../../api/models';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';

@Component({
  selector: 'app-resource-picker',
  templateUrl: './resource-picker.component.html',
  styleUrls: ['./resource-picker.component.css'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => ResourcePickerComponent),
      multi: true
    }
  ]
})
export class ResourcePickerComponent implements OnInit, OnDestroy, ControlValueAccessor {

  private sharedWithGroupSubject = new ReplaySubject<string>();
  private sharedWithGroup$ = this.sharedWithGroupSubject.asObservable().pipe(
    startWith(undefined as string | undefined),
    distinctUntilChanged(),
    shareReplay()
  );

  private subTypeSubject = new ReplaySubject<ResourceType>();
  private subType$ = this.subTypeSubject.asObservable().pipe(
    startWith(undefined as ResourceType | undefined),
    distinctUntilChanged(),
    shareReplay()
  );

  private createdBySubject = new Subject<string>();
  private createdBy$ = this.createdBySubject.asObservable().pipe(
    startWith(undefined as string | undefined),
    distinctUntilChanged(),
    shareReplay());

  // ControlValueAccessor backing field
  private propagateChangeFn: (val: string) => void;

  constructor(private backend: BackendService) {
  }

  querySubject = new Subject<string>();
  query$ = this.querySubject.asObservable().pipe(startWith(''), shareReplay(1));
  items$ = combineLatest([this.query$, this.createdBy$, this.sharedWithGroup$, this.subType$])
    .pipe(
      tap(a => console.log),
      switchMap(([q, c, g, subType]) => {
        return this.backend.searchResources(new SearchResourceRequest(q, CallType.Offer, subType, c, g, 10, 0));
      }),
      pluck('resources')
    );

  @Input()
  set createdBy(value: string | undefined) {
    this.createdBySubject.next(value);
  }

  @Input()
  set sharedWithGroup(value: string | undefined) {
    this.sharedWithGroupSubject.next(value);
  }

  @Input()
  set subType(value: string | undefined) {
    this.subTypeSubject.next(value as ResourceType);
  }

  selectedIdSubject = new Subject<string>();
  selectedId$ = this.selectedIdSubject.asObservable().pipe(
    startWith(null as string | null),
    distinctUntilChanged(),
    shareReplay()
  );

  selectedIdSub = this.selectedId$.subscribe(id => {
    this.propagateChange(id);
  });

  selectedSubject = new Subject<Resource | null>();
  selected$ = this.selectedSubject.asObservable().pipe(
    tap(a => this.propagateChange(a?.resourceId))
  );

  selectedSub = combineLatest([this.selected$, this.createdBy$]).subscribe(([selected, createdBy]) => {
    if (selected && createdBy && selected.createdBy !== createdBy) {
      this.selectedIdSubject.next(null);
    }
  });

  // This is the control value
  flexibleSubject = new Subject<string | Resource | null>();
  flexibleSub = this.flexibleSubject.asObservable()
    .pipe(
      startWith(null as string | Resource | null),
      distinctUntilChanged(),
      switchMap(res => {
        if (res === null || res === '' || res === undefined) {
          return of(null as Resource | null);
        }
        if (isResource(res)) {
          return of(res);
        }
        console.log(res)
        return this.backend.getResource(res).pipe(pluck('resource'));
      }),
      shareReplay()
    ).subscribe(res => {
      this.selectedIdSubject.next(res ? res.resourceId : null);
      this.selectedSubject.next(res);
    });

  // Begin ControlValueAccessor implementation

  propagateChange(val: string) {
    if (this.propagateChangeFn) {
      this.propagateChangeFn(val);
    }
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
    this.flexibleSubject.next(obj);
  }

  // End ControlValueAccessor implementation

  ngOnDestroy(): void {
    if (this.selectedIdSub) {
      this.selectedIdSub.unsubscribe();
    }
    if (this.flexibleSub) {
      this.flexibleSub.unsubscribe();
    }
    if (this.selectedSub) {
      this.selectedSub.unsubscribe();
    }
  }
}

function isResource(res: string | Resource): res is Resource {
  return res && (res as Resource).resourceId !== undefined;
}

function isResourceId(res: string | Resource | undefined): res is string {
  return typeof res === 'string';
}
