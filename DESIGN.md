# Các thuật ngữ

Một hệ thống chat bao gồm 3 models chính

+ user
+ room
+ event

Một user có thể join nhiều room khác nhau. 
Một room có thể có từ 2-N users.
event có thể là một message được tạo ra từ user, hoặc là một sư kiện do user sinh ra, ví dụ: user mới join vào room, user left ra khỏi room, user typing. 
event là sẽ được broadcast tới tất cả các user đang online trong hệ thống. Nếu user không online, thì các message event sẽ được lưu lại trong DB, để khi user đăng nhập lại, user có thể đọc được


# Kiến trúc của Sendpi

+ Mỗi một room sẽ được đại diện bởi một channel. Channel có thể có nhiều subscriber. Mỗi khi có một event được tạo ra trong Channel, channel sẽ dispatch event tới tất cả các subscriber.

+ Channel sẽ tự động bị xóa khỏi hệ thống sau N giây nếu như nó không có bất cứ subscribe nào. Channel sẽ được tạo lại khi có một Event mới được tạo ra từ trong Channel này. Điều này đảm bảo cho hệ thống chỉ cần giữ lại các Active Channel mà thôi. N có thể thay đổi với từng channel khác nhau, tùy theo mức độ thường xuyên của các Event của Channel.

+ Khi user connect vào hệ thống, user sẽ tự động subscribe vào một channel đại diện cho user đó. Điều này giúp cho user có thể đăng nhập trên nhiều thiết bị, và tất cả các thiết bị điều nhận được sự thay đổi từ user. 

+ Để đảm bảo event được dispatch tới tất cả các subscriber, Channel cần lưu trữ 2 danh sách:
  
   * L1: danh sách các user của room đang online -> đây chính là danh sách các subscriber của Channel
   * L2: danh sách các user của Room đang offline. Khi có message mới, các user này sẽ nhận được push notification từ hệ thống

+ Channel được tạo ra từ 2 tình huống:

  * Channel được tạo ra lần đầu tiên: Ví dụ user tạo ra Room mới. L1 = [user tạo ra room], L2 = []
  * Channel được hồi sinh lại khi có một user tạo ra event trong Channel. Lúc này, hệ thống sẽ tự động tính lại L1 của Channel, bằng cách loop qua các user trong Room, kiểm tra xem L1 của Channel của user là có rỗng hay không. L1 của User rỗng nghĩa là user đang offline, ngược lại là user online. 

+ Làm sao một user biết nên connect tới Channel nào khi user kết nối với hệ thống? Mỗi khi user kết nối tới hệ thống, hệ thống sẽ lấy ra tất cả các Room mà user tham gia, các Room này được sắp xếp theo thứ tự thời gian của message cuối cùng giảm dần. Nếu user tham gia 2 rooms Ri và Rj. Nếu message cuối cùng của Ri > Rj, thì Channel của Rj sẽ không thể tồn tại nếu như Channel của Ri không tồn tại, điều này giúp cho việc tìm kiếm các Channel mà user nên tham gia vào trở nên rất đơn giản. Chúng ta chỉ cần tìm room lớn nhất mà không tồn tại Channel ứng với nó trong hệ thống.

# Scale hệ thống:

Chúng ta đã phân hệ thống thành 2 tầng: 

Tầng user - là các connection ứng với một session của một user. 
Tầng Channel

Chúng ta tạm gọi các mỗi một session của User hoặc một channel à 1 Process (có thể là goroutine của Go, hoặc là GenServer của Elixir)

Chúng ta có 1 cụm server cho các Process của User, và một cụm server cho các 
Process của Channel. 
Trong các cụm này, mỗi Process sẽ được shared bởi ID của Process (ví dụ ID của user, hoặc ID của Room)