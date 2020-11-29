import {MarkdownRenderer} from './renderer';
import {async, inject, TestBed} from '@angular/core/testing';
import {CreateOfferComponent} from '../offers/create-offer/create-offer.component';
import {CommonModule} from '@angular/common';
import {BrowserModule, DomSanitizer} from '@angular/platform-browser';

describe('renderer', () => {

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [BrowserModule],
    }).compileComponents();
  }));

  it('create an instance', inject([DomSanitizer], (domSanitizer: DomSanitizer) => {
    const renderer = new MarkdownRenderer(domSanitizer);
    const result = renderer.text('<commonpool-user>Hello</commonpool-user>');
    console.log(result);
    expect(renderer).toBeTruthy();
  }));
});
