import {Component, Input, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {AcceptOfferRequest, DeclineOfferRequest, GetOffersRequest, Offer} from '../../api/models';

@Component({
  selector: 'app-offer-list',
  templateUrl: './offer-list.component.html',
  styleUrls: ['./offer-list.component.css']
})
export class OfferListComponent implements OnInit {

  constructor(private backend: BackendService) {

  }

  offers: Offer[] = [];

  private _userId: string;
  @Input()
  set userId(val: string) {
    this._userId = val;
    this.refresh();
  }

  get userId() {
    return this._userId;
  }

  ngOnInit(): void {
  }

  accept(id: string) {
    this.backend.acceptOffer(new AcceptOfferRequest(id)).subscribe(res => {
      console.log(res);
      this.refresh();
    });
  }

  decline(id: string) {
    this.backend.declineOffer(new DeclineOfferRequest(id)).subscribe(res => {
      console.log(res);
      this.refresh();
    });
  }

  refresh() {
    if (!this.userId) {
      this.offers = [];
      return;
    }
    this.backend.getOffers(new GetOffersRequest()).subscribe(offers => {
      console.log(offers);
      this.offers = offers.offers;
    });
  }

}
