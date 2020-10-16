import {async, ComponentFixture, fakeAsync, flush, flushMicrotasks, TestBed, tick} from '@angular/core/testing';

import {ResourceListViewComponent} from './resource-list-view.component';
import {BackendService} from '../../api/backend.service';
import {Resource, ResourceType, SearchResourcesResponse} from '../../api/models';
import {FormsModule} from '@angular/forms';
import {By} from '@angular/platform-browser';
import SpyObj = jasmine.SpyObj;
import {TestScheduler} from 'rxjs/testing';
import {defer, of} from 'rxjs';
import {delay, shareReplay, take, tap} from 'rxjs/operators';

describe('ResourceListViewComponent', () => {
  let component: ResourceListViewComponent;
  let fixture: ComponentFixture<ResourceListViewComponent>;
  let backend: SpyObj<BackendService>;
  let scheduler: TestScheduler;

  beforeEach(async(() => {
    scheduler = new TestScheduler((actual, expected) => {
      expect(actual).toEqual(expected);
    });
    backend = jasmine.createSpyObj<BackendService>(['searchResources']);
    TestBed.configureTestingModule({
      imports: [FormsModule],
      declarations: [ResourceListViewComponent],
      providers: [{
        provide: BackendService,
        useValue: backend
      }]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceListViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display a search box', () => {
    const input = fixture.debugElement.query(By.css('#searchBox'));
    expect(input).toBeTruthy();
  });


  it('should make a request when user inputs a query', fakeAsync(() => {
    const resources = new SearchResourcesResponse([
      new Resource(
        'eb0f735d-7d00-466b-84bb-2e6d0f83aa9e',
        'summary',
        'description',
        ResourceType.Offer,
        10,
        20,
        30
      )], 1, 10, 0);

    const resources$ = defer(() => of(resources).pipe(delay(100)));
    const spy = backend.searchResources.and.returnValue(resources$);

    expect(spy).toHaveBeenCalledTimes(0);
    assertHasResourceCount(0);

    enterSearchTerm('new input');

    tick(150);
    fixture.detectChanges();

    assertHasPendingIndicator(true);
    assertHasResourceCount(0);
    expect(spy).toHaveBeenCalledTimes(1);

    tick(500);
    fixture.detectChanges();

    assertHasPendingIndicator(false);
    assertHasResourceCount(1);

  }));

  function assertHasResourceCount(numberOfResources: number) {
    expect(fixture.debugElement.queryAll(By.css('.resource')).length).toBe(numberOfResources);
  }

  function assertHasPendingIndicator(toBePresent = true) {
    const actual = fixture.debugElement.query(By.css('#pendingIndicator'));
    if (toBePresent) {
      expect(actual).toBeTruthy();
    } else {
      expect(actual).toBeFalsy();
    }
  }

  function enterSearchTerm(searchTerm: string) {
    const inputEl = fixture.debugElement.query(By.css('#searchBox')).nativeElement as HTMLInputElement;
    inputEl.value = searchTerm;
    inputEl.dispatchEvent(new Event('input'));
  }
});

