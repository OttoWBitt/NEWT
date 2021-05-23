import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ArtifactDialogComponent } from './artifact-dialog.component';

describe('ArtifactDialogComponent', () => {
  let component: ArtifactDialogComponent;
  let fixture: ComponentFixture<ArtifactDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ArtifactDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ArtifactDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
