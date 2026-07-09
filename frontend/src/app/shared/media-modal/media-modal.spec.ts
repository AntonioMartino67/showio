import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MediaModal } from './media-modal';

describe('MediaModal', () => {
  let component: MediaModal;
  let fixture: ComponentFixture<MediaModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MediaModal],
    }).compileComponents();

    fixture = TestBed.createComponent(MediaModal);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
