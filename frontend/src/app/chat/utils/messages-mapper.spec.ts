import {Message, MessageSubType, MessageType} from '../../api/models';
import {MapMessages} from './messages-mapper';

function ADate(daysFromNow: number, hour: number): Date {
  const dt = new Date();
  dt.setDate(dt.getDate() + daysFromNow);
  dt.setHours(hour);
  return dt;
}

function AMessage(id: string, sentBy: string, sentAt: Date): Message {
  return {
    sentAtDate: sentAt,
    sentAtDistance: '',
    sentAtHour: '',
    id,
    channelId: '',
    text: '',
    visibleToUser: '',
    sentById: sentBy,
    sentByUsername: '',
    attachments: [],
    blocks: [],
    sentAt: '',
    messageSubType: MessageSubType.UserMessage,
    messageType: MessageType.NormalMessage
  };
}

describe('MessagesMapper', () => {
  it('should group messages by date ', () => {

    const message1 = AMessage('1', 'user1', ADate(1, 10));
    const message2 = AMessage('2', 'user1', ADate(1, 11));
    const message3 = AMessage('3', 'user1', ADate(1, 12));
    const message4 = AMessage('4', 'user1', ADate(2, 10));

    const messages: Message[] = [message1, message2, message3, message4];
    const mappedMessages = MapMessages(messages);

    expect(mappedMessages.length).toBe(2);
    expect(mappedMessages[0].messages.length).toBe(3);
    expect(mappedMessages[0].messages[0].id).toBe('1');
    expect(mappedMessages[0].messages[1].id).toBe('2');
    expect(mappedMessages[0].messages[2].id).toBe('3');
    expect(mappedMessages[1].messages.length).toBe(1);
    expect(mappedMessages[1].messages[0].id).toBe('4');

  });

  it('should group messages by user ', () => {

    const message1 = AMessage('1', 'user1', ADate(1, 10));
    const message2 = AMessage('2', 'user1', ADate(1, 10));
    const message3 = AMessage('3', 'user1', ADate(1, 10));
    const message4 = AMessage('4', 'user2', ADate(1, 10));
    const message5 = AMessage('5', 'user2', ADate(1, 10));

    const messages: Message[] = [message1, message2, message3, message4, message5];
    const mappedMessages = MapMessages(messages);

    expect(mappedMessages.length).toBe(2);
    expect(mappedMessages[0].messages.length).toBe(3);
    expect(mappedMessages[0].messages[0].sentById).toBe('user1');
    expect(mappedMessages[0].messages[1].sentById).toBe('user1');
    expect(mappedMessages[0].messages[2].sentById).toBe('user1');
    expect(mappedMessages[1].messages.length).toBe(2);
    expect(mappedMessages[1].messages[0].sentById).toBe('user2');
    expect(mappedMessages[1].messages[0].sentById).toBe('user2');

  });

  it('should group messages by both user and date ', () => {

    const message1 = AMessage('1', 'user1', ADate(1, 10));
    const message2 = AMessage('2', 'user1', ADate(1, 11));
    const message3 = AMessage('3', 'user2', ADate(1, 12));
    const message4 = AMessage('4', 'user2', ADate(2, 13));
    const message5 = AMessage('5', 'user1', ADate(2, 14));
    const message6 = AMessage('6', 'user2', ADate(2, 15));
    const message7 = AMessage('7', 'user2', ADate(2, 16));

    const messages: Message[] = [message1, message2, message3, message4, message5, message6, message7];
    const mappedMessages = MapMessages(messages);

    expect(mappedMessages.length).toBe(5);

    expect(mappedMessages[0].messages.length).toBe(2);
    expect(mappedMessages[0].messages[0].sentById).toBe('user1');
    expect(mappedMessages[0].messages[1].sentById).toBe('user1');

    expect(mappedMessages[1].messages.length).toBe(1);
    expect(mappedMessages[1].messages[0].sentById).toBe('user2');

    expect(mappedMessages[2].messages.length).toBe(1);
    expect(mappedMessages[2].messages[0].sentById).toBe('user2');

    expect(mappedMessages[3].messages.length).toBe(1);
    expect(mappedMessages[3].messages[0].sentById).toBe('user1');

    expect(mappedMessages[4].messages.length).toBe(2);
    expect(mappedMessages[4].messages[0].sentById).toBe('user2');
    expect(mappedMessages[4].messages[0].sentById).toBe('user2');

  });
});
