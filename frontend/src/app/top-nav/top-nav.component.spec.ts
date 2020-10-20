import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {TopNavComponent} from './top-nav.component';
import {AuthService} from '../auth.service';


describe('TopNavComponent', () => {
  let component: TopNavComponent;
  let fixture: ComponentFixture<TopNavComponent>;

  const mockAuthService = {};

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TopNavComponent],
      providers: [
        {
          provide: AuthService,
          useValue: mockAuthService
        }
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TopNavComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
