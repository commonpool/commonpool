import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupLinkComponent } from './group-link.component';

describe('GroupLinkComponent', () => {
  let component: GroupLinkComponent;
  let fixture: ComponentFixture<GroupLinkComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GroupLinkComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupLinkComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
