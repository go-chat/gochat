# Các thuật ngữ

Một hệ thống chat bao gồm 3 models chính

+ User
+ Room
+ Event

Một User có thể join nhiều Room khác nhau. 
Một Room có thể có từ 1-N Users.
Event có thể là một message được tạo ra từ User, hoặc là một sư kiện do User sinh ra, ví dụ: User mới join vào Room, User left ra khỏi Room, User typing. 
Event là sẽ được broadcast tới tất cả các User của Room, đang online trong hệ thống. Nếu User không online, thì các message Event sẽ được lưu lại trong DB, để khi User đăng nhập lại, User có thể đọc được


# Kiến trúc của Sendpi

+ Mỗi một Room sẽ được đại diện bởi một Channel. Channel có thể có nhiều subscriber. Mỗi khi có một Event được tạo ra trong Channel, Channel sẽ dispatch Event tới tất cả các subscriber.

+ Channel sẽ tự động bị xóa khỏi hệ thống sau N giây nếu như nó không có bất cứ subscribe nào. Channel sẽ được tạo lại khi có một Event mới được tạo ra từ trong Channel này. Điều này đảm bảo cho hệ thống chỉ cần giữ lại các Active Channel mà thôi. N có thể thay đổi với từng Channel khác nhau, tùy theo mức độ thường xuyên của các Event của Channel.

+ Khi User connect vào hệ thống, User sẽ tự động subscribe vào một Channel đại diện cho User đó. Điều này giúp cho User có thể đăng nhập trên nhiều thiết bị, và tất cả các thiết bị điều nhận được sự thay đổi từ User. 

+ Để đảm bảo Event được dispatch tới tất cả các subscriber, Channel cần lưu trữ 2 danh sách:
  
   * L1: danh sách các User của Room đang online -> đây chính là danh sách các subscriber của Channel
   * L2: danh sách các User của Room đang offline. Khi có message mới, các User này sẽ nhận được push notification từ hệ thống

+ Channel được tạo ra từ 2 tình huống:

  * Channel được tạo ra lần đầu tiên: Ví dụ User tạo ra Room mới. L1 = [User tạo ra Room], L2 = []
  * Channel được tái kích hoạt lại khi có một User Event mới trong Channel. Lúc này, hệ thống sẽ tự động tính lại L1 của Channel, bằng cách loop qua các User trong Room, kiểm tra xem L1 của Channel của User là có rỗng hay không. L1 của User rỗng nghĩa là User đang offline, ngược lại là User online. 

+ Làm sao một User biết nên connect tới Channel nào khi User kết nối với hệ thống? Mỗi khi User kết nối tới hệ thống, hệ thống sẽ lấy ra tất cả các Room mà User tham gia, các Room này được sắp xếp theo thứ tự thời gian của message cuối cùng giảm dần. Nếu User tham gia 2 Rooms Ri và Rj. Nếu message cuối cùng của Ri > Rj, thì Channel của Rj sẽ không thể tồn tại nếu như Channel của Ri không tồn tại, điều này giúp cho việc tìm kiếm các Channel mà User nên tham gia vào trở nên rất đơn giản. Chúng ta chỉ cần tìm Room lớn nhất mà không tồn tại Channel ứng với nó trong hệ thống.

# Scale hệ thống:

Chúng ta đã phân hệ thống thành 2 tầng: 

Tầng User - là các connection ứng với một session của một User. 
Tầng Channel

Chúng ta tạm gọi các mỗi một session của User hoặc một Channel à 1 Process (có thể là goroutine của Go, hoặc là GenServer của Elixir)

Chúng ta có 1 cụm server cho các Process của User, và một cụm server cho các 
Process của Channel. 
Trong các cụm này, mỗi Process sẽ được shard bởi ID của Process (ví dụ ID của User, hoặc ID của Room)