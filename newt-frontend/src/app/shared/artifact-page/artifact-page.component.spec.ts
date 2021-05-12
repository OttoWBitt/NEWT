import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ArtifactPageComponent } from './artifact-page.component';

describe('ArtifactPageComponent', () => {
  let component: ArtifactPageComponent;
  let fixture: ComponentFixture<ArtifactPageComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ArtifactPageComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ArtifactPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
