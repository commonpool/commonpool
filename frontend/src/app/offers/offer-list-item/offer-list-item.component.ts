import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {
  Action,
  ConfirmBorrowedResourceReturned,
  ConfirmResourceBorrowed,
  ConfirmResourceTransferred,
  ConfirmServiceProvidedRequest,
  Offer
} from '../../api/models';
import {combineLatest, of, ReplaySubject, Subject} from 'rxjs';
import {map, tap} from 'rxjs/operators';
import {AuthService} from '../../auth.service';
import {BackendService} from '../../api/backend.service';

@Component({
  selector: 'app-offer-list-item',
  templateUrl: './offer-list-item.component.html',
  styleUrls: ['./offer-list-item.component.css']
})
export class OfferListItemComponent implements OnInit {

  constructor(public auth: AuthService, private backend: BackendService) {
  }

  private offerSubject = new ReplaySubject<Offer>(1);
  public offerSubject$ = this.offerSubject.asObservable();
  public allOfferActions$ = this.offerSubject$.pipe(
    map(o => o.actions)
  );

  public offerItemActions$ = this.allOfferActions$.pipe(
    map((o) => {
      const result: { [key: string]: Action[] } = {};
      for (const action of o) {
        if (action.offerItemId !== '') {
          if (!result[action.offerItemId]) {
            result[action.offerItemId] = [];
          }
          result[action.offerItemId].push(action);
        }
      }
      return result;
    })
  );
  public offerActions$ = this.allOfferActions$.pipe(
    map((o) => {
      return o.filter(action => !action.offerItemId);
    })
  );

  @Input()
  set offer(value: Offer) {
    this.offerSubject.next(value);

  }

  @Output()
  refresh = new EventEmitter();

  ngOnInit(): void {
  }

}

