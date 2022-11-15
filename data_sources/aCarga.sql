CREATE PROCEDURE RegistrarMetodoPago -holaaaa

@CodMetodoPago VARCHAR(7), 
*
ESTO ES UN ERROR
*/
@Descripcion VARCHAR(20),
@Premiacion BIT
AS
BEGIN
SELECT * FROM MetodosPago
INSERT INTO MetodosPago(CodMetodoPago, Descripcion, Premiacion) VALUES(@CodMetodoPago, @Descripcion, @Premiacion)
END
GO
