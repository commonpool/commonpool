import {Component, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {GetTradingHistoryRequest} from '../../api/models';

@Component({
  selector: 'app-trading-history',
  templateUrl: './trading-history.component.html',
  styleUrls: ['./trading-history.component.css']
})
export class TradingHistoryComponent implements OnInit {

  constructor(private backend: BackendService) {
  }

  history$ = this.backend.getTradingHistory(new GetTradingHistoryRequest([]));

  ngOnInit(): void {
  }

}
