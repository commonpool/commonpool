import {Component, forwardRef, Input, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {combineLatest, of, ReplaySubject, Subject, Subscription} from 'rxjs';
import {distinctUntilChanged, map, pluck, shareReplay, startWith, switchMap, tap} from 'rxjs/operators';
import {Resource, ResourceType, SearchResourceRequest} from '../../api/models';
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
export class ResourcePickerComponent implements OnInit, ControlValueAccessor {


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
  query$ = this.querySubject.asObservable().pipe(startWith(''));
  items$ = combineLatest([this.query$, this.createdBy$])
    .pipe(
      switchMap(([q, c]) => this.backend.searchResources(new SearchResourceRequest(q, ResourceType.Offer, c, 10, 0))),
      pluck('resources')
    );

  @Input()
  set createdBy(value: string | undefined) {
    console.log('Created by : ', value);
    this.createdBySubject.next(value);
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
    tap(a => console.log('selected', a)),
    tap(a => this.propagateChange(a?.id))
  );

  selectedSub = combineLatest([this.selected$, this.createdBy$]).subscribe(([selected, createdBy]) => {
    if (selected && createdBy && selected.createdById !== createdBy) {
      console.log('selected', selected, 'createdBy', createdBy);
      this.selectedIdSubject.next(null);
    }
  });

  // This is the control value
  flexibleSubject = new Subject<string | Resource | null>();
  flexible$ = this.flexibleSubject.asObservable()
    .pipe(
      startWith(null as string | Resource | null),
      distinctUntilChanged(),
      switchMap(res => {
        console.log('hub', res);
        if (res === null) {
          return of(null as Resource | null);
        }
        if (isResource(res)) {
          return of(res);
        }
        return this.backend.getResource(res).pipe(pluck('resource'));
      }),
      shareReplay()
    ).subscribe(res => {
      this.selectedIdSubject.next(res ? res.id : null);
      this.selectedSubject.next(res);
    });


  // Begin ControlValueAccessor implementation

  propagateChange(val: string) {
    console.log('propagate', val);
    if (this.propagateChangeFn) {
      console.log('propagating', val);
      this.propagateChangeFn(val);
    }
  }


  ngOnInit(): void {
    this.querySubject.asObservable().subscribe(q => console.log(q));
  }

  registerOnChange(fn: any): void {
    this.propagateChangeFn = fn;
  }

  registerOnTouched(fn: any): void {
  }

  setDisabledState(isDisabled: boolean): void {
  }

  writeValue(obj: any): void {
    console.log('write', obj);
    this.flexibleSubject.next(obj);
  }

  // End ControlValueAccessor implementation
}


function isResource(res: string | Resource): res is Resource {
  return res && (res as Resource).id !== undefined;
}

function isResourceId(res: string | Resource | undefined): res is string {
  return typeof res === 'string';
}
